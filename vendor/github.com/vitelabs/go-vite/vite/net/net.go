package net

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/crypto/ed25519"
	"github.com/vitelabs/go-vite/log15"
	"github.com/vitelabs/go-vite/p2p"
	"github.com/vitelabs/go-vite/p2p/vnode"
	"github.com/vitelabs/go-vite/vite/net/protos"
)

var netLog = log15.New("module", "net")

var errNetIsRunning = errors.New("network is already running")
var errNetIsNotRunning = errors.New("network is not running")

type Config struct {
	Single             bool
	FileListenAddress  string
	TraceEnabled       bool
	ForwardStrategy    string   // default `cross`
	AccessControl      string   `json:"AccessControl"` // producer special any
	AccessAllowKeys    []string `json:"AccessAllowKeys"`
	AccessDenyKeys     []string `json:"AccessDenyKeys"`
	BlackBlockHashList []string

	MinePrivateKey ed25519.PrivateKey
	P2PPrivateKey  ed25519.PrivateKey
	Chain
	Verifier
}

const DefaultForwardStrategy = "cross"
const DefaultFilePort = 8484
const ID = 1
const maxNeighbors = 100

type net struct {
	Config
	nodeID   vnode.NodeID
	peers    *peerSet
	*syncer  // use pointer but not interface, because syncer can be start/stop, but interface has no start/stop method
	*fetcher // use pointer but not interface, because fetcher can be start/stop, but interface has no start/stop method
	*broadcaster
	reader     *cacheReader
	downloader syncDownloader
	BlockSubscriber
	server    *syncServer
	handlers  *msgHandlers
	query     *queryHandler
	hb        *heartBeater
	running   int32
	log       log15.Logger
	tracer    Tracer
	sn        *sbpn
	consensus Consensus
	p2p       p2p.P2P
}

func (n *net) ProtoData() (height uint64, head types.Hash, genesis types.Hash) {
	genesis = n.Chain.GetGenesisSnapshotBlock().Hash
	current := n.Chain.GetLatestSnapshotBlock()

	height = current.Height
	head = current.Hash

	return
}

func (n *net) ReceiveHandshake(msg *p2p.HandshakeMsg) (level p2p.Level, err error) {
	var id = msg.ID.String()
	var key string

	if msg.Key != nil {
		addr := types.PubkeyToAddress(msg.Key)
		if n.sn != nil && n.sn.isSbp(addr) {
			level = p2p.Superior
		}
	}

	for _, key2 := range n.AccessDenyKeys {
		if key2 == id || key2 == key {
			err = p2p.PeerNoPermission
			return
		}
	}

	if n.Config.AccessControl == "any" {
		return
	}

	if strings.Contains(n.Config.AccessControl, "producer") && level == p2p.Superior {
		return
	}

	for _, key2 := range n.AccessAllowKeys {
		if key2 == id || key2 == key {
			return
		}
	}

	err = p2p.PeerNoPermission
	return
}

func (n *net) Handle(msg p2p.Msg) error {
	p := n.peers.get(msg.Sender.ID())
	if p != nil {
		return n.handlers.handle(msg, p)
	} else {
		return errPeerNotExist
	}
}

func (n *net) SetState(state []byte, peer p2p.Peer) {
	var heartbeat = &protos.State{}

	err := proto.Unmarshal(state, heartbeat)
	if err != nil {
		n.log.Error(fmt.Sprintf("failed to unmarshal heartbeat message from %s: %v", peer, err))
		return
	}

	p := n.peers.get(peer.ID())
	if p != nil {
		var head types.Hash
		head, err = types.BytesToHash(heartbeat.Head)
		if err != nil {
			return
		}

		p.SetHead(head, heartbeat.Height)

		// max 100 neighbors
		var count = len(heartbeat.Peers)
		if count > maxNeighbors {
			count = maxNeighbors
		}
		var pl = make([]peerConn, count)
		for i := 0; i < count; i++ {
			hp := heartbeat.Peers[i]
			pl[i] = peerConn{
				id:  hp.ID,
				add: hp.Status != protos.State_Disconnected,
			}
		}
		p.setPeers(pl, heartbeat.Patch)
	}
}

func (n *net) OnPeerAdded(peer p2p.Peer) (err error) {
	p := newPeer(peer, netLog.New("peer", peer.ID()))

	err = n.peers.add(p)
	if err != nil {
		return
	}

	go n.syncer.start()

	return nil
}

func (n *net) OnPeerRemoved(peer p2p.Peer) (err error) {
	_, err = n.peers.remove(peer.ID())

	return
}

func New(cfg Config) Net {
	// for test
	if cfg.Single {
		return mock(cfg)
	}

	cfg.BlackBlockHashList = append(cfg.BlackBlockHashList, []string{
		"6771bc124fed97302328c13fb9a97919c8963b7b1f79a431091c7ace00ec28f4",
		"3963c532b43d476f1cadd01dc36cd5e157b40c86f6848665549e5959626efd39",
		"f8a9579e36d605e0f1d9e3d2de96e798d3a8218d771cc4d996cf304fac92ed40",
		"3cdbdd9777eecdd1238675cd0e25b94742af6d5603a926cd331a3b9ce07a1f73",
		"059aee3dcb367197b2599bcf40b1b02e5bd763cd6c6c65edc5fbc90807129ed3",
		"3fccebb35716e448a2c08341d0bd5a30b8122a7dd3391eff081013931f25f922",
	}...)

	peers := newPeerSet()

	feed := newBlockFeeder()
	feed.setBlackHashList(cfg.BlackBlockHashList)

	forward := createForardStrategy(cfg.ForwardStrategy, peers)
	broadcaster := newBroadcaster(peers, cfg.Verifier, feed, newMemBlockStore(1000), forward, nil, cfg.Chain)

	receiver := &safeBlockNotifier{
		blockFeeder: feed,
		Verifier:    cfg.Verifier,
	}

	var id peerId
	id, _ = vnode.Bytes2NodeID(cfg.P2PPrivateKey.PubByte())
	syncConnFac := &defaultSyncConnectionFactory{
		chain:   cfg.Chain,
		peers:   peers,
		id:      id,
		peerKey: cfg.P2PPrivateKey,
		mineKey: cfg.MinePrivateKey,
	}
	downloader := newExecutor(50, 10, peers, syncConnFac)
	syncServer := newSyncServer(cfg.FileListenAddress, cfg.Chain, syncConnFac)

	reader := newCacheReader(cfg.Chain, cfg.Verifier, downloader)
	reader.setBlackHashList(cfg.BlackBlockHashList)

	syncer := newSyncer(cfg.Chain, peers, reader, downloader, 10*time.Minute)
	syncer.setBlackHashList(cfg.BlackBlockHashList)

	fetcher := newFetcher(peers, receiver)
	fetcher.setBlackHashList(cfg.BlackBlockHashList)

	syncer.SubscribeSyncStatus(fetcher.subSyncState)
	syncer.SubscribeSyncStatus(broadcaster.subSyncState)

	n := &net{
		Config:          cfg,
		BlockSubscriber: feed,
		peers:           peers,
		syncer:          syncer,
		reader:          reader,
		fetcher:         fetcher,
		broadcaster:     broadcaster,
		downloader:      downloader,
		server:          syncServer,
		handlers:        newHandlers("vite"),
		hb:              newHeartBeater(peers, cfg.Chain),
		log:             netLog,
	}

	var err error
	err = n.handlers.register(&stateHandler{
		maxNeighbors: 100,
		peers:        peers,
	})
	if err != nil {
		panic(fmt.Errorf("cannot register handler: broadcaster: %v", err))
	}

	n.query, err = newQueryHandler(cfg.Chain)
	if err != nil {
		panic(fmt.Errorf("cannot construct query handler: %v", err))
	}

	// GetSubLedgerCode, CodeGetSnapshotBlocks, CodeGetAccountBlocks, GetChunkCode
	if err = n.handlers.register(n.query); err != nil {
		panic(fmt.Errorf("cannot register handler: query: %v", err))
	}

	// CodeNewSnapshotBlock, CodeNewAccountBlock
	if err = n.handlers.register(broadcaster); err != nil {
		panic(fmt.Errorf("cannot register handler: broadcaster: %v", err))
	}

	// CodeSnapshotBlocks, CodeAccountBlocks
	if err = n.handlers.register(fetcher); err != nil {
		panic(fmt.Errorf("cannot register handler: fetcher: %v", err))
	}

	// CodeSnapshotBlocks, CodeAccountBlocks
	if err = n.handlers.register(syncer); err != nil {
		panic(fmt.Errorf("cannot register handler: syncer: %v", err))
	}

	// trace
	if cfg.TraceEnabled {
		var p2pPub = cfg.P2PPrivateKey.PubByte()
		n.tracer = newTracer(hex.EncodeToString(p2pPub), peers, forward)
	} else {
		n.tracer = newMockTracer()
	}
	if err = n.handlers.register(n.tracer); err != nil {
		panic(fmt.Errorf("cannot register handler: tracer: %v", err))
	}

	return n
}

type heartBeater struct {
	chain     chainReader
	last      time.Time
	lastPeers map[peerId]struct{}
	ps        *peerSet
}

func newHeartBeater(ps *peerSet, chain chainReader) *heartBeater {
	return &heartBeater{
		chain:     chain,
		lastPeers: make(map[peerId]struct{}),
		ps:        ps,
	}
}

func (h *heartBeater) state() []byte {
	current := h.chain.GetLatestSnapshotBlock()

	var heartBeat = &protos.State{
		Peers:     nil,
		Patch:     true,
		Head:      current.Hash.Bytes(),
		Height:    current.Height,
		Timestamp: time.Now().Unix(),
	}

	idMap := h.ps.idMap()

	if time.Now().Sub(h.last) > time.Minute {
		heartBeat.Patch = false
		h.lastPeers = make(map[vnode.NodeID]struct{})
	}

	var id vnode.NodeID
	var ok bool
	for id = range h.lastPeers {
		if _, ok = idMap[id]; ok {
			delete(idMap, id)
			continue
		}
		heartBeat.Peers = append(heartBeat.Peers, &protos.State_Peer{
			ID:     id.Bytes(),
			Status: protos.State_Disconnected,
		})
	}

	for id = range idMap {
		heartBeat.Peers = append(heartBeat.Peers, &protos.State_Peer{
			ID:     id.Bytes(),
			Status: protos.State_Connected,
		})
	}

	data, err := proto.Marshal(heartBeat)
	if err != nil {
		return nil
	}

	h.lastPeers = idMap

	return data
}

func (n *net) Init(consensus Consensus, reader IrreversibleReader) {
	n.consensus = consensus
	n.syncer.irreader = reader
	n.reader.irreader = reader
}

func (n *net) State() []byte {
	return n.hb.state()
}

func (n *net) Start(svr p2p.P2P) (err error) {
	if atomic.CompareAndSwapInt32(&n.running, 0, 1) {
		n.nodeID = svr.Config().Node().ID
		n.p2p = svr

		if err = n.server.start(); err != nil {
			return
		}

		n.downloader.start()

		n.reader.start()

		n.query.start()

		n.fetcher.start()

		var addr types.Address
		if len(n.MinePrivateKey) != 0 {
			addr = types.PubkeyToAddress(n.MinePrivateKey.PubByte())
			setNodeExt(n.MinePrivateKey, svr.Node())
			n.fetcher.setSBP()
		}
		n.sn = newSbpn(addr, n.peers, svr, n.consensus)
		if discv := svr.Discovery(); discv != nil {
			discv.SubscribeNode(n.sn.receiveNode)
		}

		return
	}

	return errNetIsRunning
}

func (n *net) Stop() error {
	if atomic.CompareAndSwapInt32(&n.running, 1, 0) {
		n.reader.stop()

		n.syncer.stop()

		n.downloader.stop()

		_ = n.server.stop()

		n.query.stop()

		n.fetcher.stop()

		n.sn.clean()

		return nil
	}

	return errNetIsNotRunning
}

func (n *net) Trace() {
	if n.tracer != nil {
		n.tracer.Trace()
	}
}

func (n *net) Info() NodeInfo {
	info := NodeInfo{
		NodeInfo: n.p2p.Info(),
		Latency:  n.broadcaster.Statistic(),
		Peers:    n.peers.info(),
		Height:   n.Chain.GetLatestSnapshotBlock().Height,
	}

	if n.server != nil {
		info.Server = n.server.status()
	}

	return info
}

type NodeInfo struct {
	p2p.NodeInfo
	Height  uint64           `json:"height"`
	Peers   []PeerInfo       `json:"peers"`
	Latency []int64          `json:"latency"` // [0,1,12,24]
	Server  FileServerStatus `json:"server"`
}
