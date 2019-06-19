package p2p

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/vitelabs/go-vite/common/types"

	"github.com/vitelabs/go-vite/p2p/vnode"
)

type mockProtocol struct {
	mu       sync.Mutex
	peers    map[vnode.NodeID]Peer
	errFac   func() error
	interval time.Duration
}

func (m *mockProtocol) ProtoData() (height uint64, head types.Hash, genesis types.Hash) {
	return 0, types.Hash{1, 2, 3}, types.Hash{4, 5, 6}
}

func (m *mockProtocol) Name() string {
	return "mock"
}

func (m *mockProtocol) ReceiveHandshake(msg *HandshakeMsg) (level Level, err error) {
	return
}

func (m *mockProtocol) Handle(msg Msg) error {
	fmt.Printf("receive message from %s code: %d, id: %d, length: %d\n", msg.Sender.Address(), msg.Code, msg.Id, len(msg.Payload))

	if m.errFac == nil {
		return nil
	}

	return m.errFac()
}

func (m *mockProtocol) State() []byte {
	return nil
}

func (m *mockProtocol) SetState(state []byte, peer Peer) {
	return
}

func (m *mockProtocol) OnPeerAdded(peer Peer) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.peers == nil {
		m.peers = make(map[vnode.NodeID]Peer)
	}

	if _, ok := m.peers[peer.ID()]; ok {
		return errors.New("peer existed")
	}

	m.peers[peer.ID()] = peer

	go func(peer Peer) {
		var i uint32

		for {
			<-time.After(m.interval)

			err := peer.WriteMsg(Msg{
				Code:    0,
				Id:      i,
				Payload: []byte("hello"),
			})

			if err != nil {
				fmt.Printf("mock protocol write message to %s error: %v", peer, err)
				return
			}

			i++
		}
	}(peer)

	return nil
}

func (m *mockProtocol) OnPeerRemoved(peer Peer) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.peers[peer.ID()]; ok {
		delete(m.peers, peer.ID())
		return nil
	}

	return errors.New("peer not existed")
}
