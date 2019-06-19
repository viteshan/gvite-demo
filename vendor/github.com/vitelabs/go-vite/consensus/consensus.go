package consensus

import (
	"sync"
	"time"

	"github.com/vitelabs/go-vite/common"
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/consensus/cdb"
	"github.com/vitelabs/go-vite/consensus/core"
	"github.com/vitelabs/go-vite/ledger"
	"github.com/vitelabs/go-vite/log15"
	"github.com/vitelabs/go-vite/pool/lock"
)

// Verifier is the interface that can verify block consensus.
type Verifier interface {
	VerifyAccountProducer(block *ledger.AccountBlock) (bool, error)
	VerifyABsProducer(abs map[types.Gid][]*ledger.AccountBlock) ([]*ledger.AccountBlock, error)
	VerifySnapshotProducer(block *ledger.SnapshotBlock) (bool, error)
}

// Event will trigger when the snapshot block needs production
type Event struct {
	Gid     types.Gid
	Address types.Address
	Stime   time.Time
	Etime   time.Time

	Timestamp         time.Time // add to block
	SnapshotTimeStamp time.Time // add to block

	VoteTime    time.Time // voteTime
	PeriodStime time.Time // start time for period
	PeriodEtime time.Time // end time for period
}

// ProducersEvent describes all SBP in one period
type ProducersEvent struct {
	Addrs []types.Address
	Index uint64
	Gid   types.Gid
}

// Subscriber provide an interface to consensus event
type Subscriber interface {
	Subscribe(gid types.Gid, id string, addr *types.Address, fn func(Event))
	UnSubscribe(gid types.Gid, id string)
	SubscribeProducers(gid types.Gid, id string, fn func(event ProducersEvent))
}

// Reader can read consensus result
type Reader interface {
	ReadByIndex(gid types.Gid, index uint64) ([]*Event, uint64, error)
	VoteTimeToIndex(gid types.Gid, t2 time.Time) (uint64, error)
	VoteIndexToTime(gid types.Gid, i uint64) (*time.Time, *time.Time, error)
}

// APIReader is just provided for RPC api
type APIReader interface {
	ReadVoteMap(t time.Time) ([]*VoteDetails, *ledger.HashHeight, error)
	ReadSuccessRate(start, end uint64) ([]map[types.Address]*cdb.Content, error)
}

// Life define the life cycle for consensus component
type Life interface {
	Start()
	Init() error
	Stop()
}

// Consensus include all interface for consensus
type Consensus interface {
	Verifier
	Subscriber
	Reader
	Life
	API() APIReader
	SBPReader() core.SBPStatReader
}

// update committee result
type consensus struct {
	common.LifecycleStatus

	mLog log15.Logger

	genesis time.Time

	rw       *chainRw
	rollback lock.ChainRollback

	snapshot  *snapshotCs
	contracts *contractsCs

	dposWrapper *dposReader

	// subscribes map[types.Gid]map[string]*subscribeEvent
	subscribes sync.Map

	api APIReader

	wg     sync.WaitGroup
	closed chan struct{}
}

func (cs *consensus) SBPReader() core.SBPStatReader {
	return cs.snapshot
}

func (cs *consensus) API() APIReader {
	return cs.api
}

// NewConsensus instantiates a new consensus object
func NewConsensus(ch Chain, rollback lock.ChainRollback) Consensus {
	log := log15.New("module", "consensus")
	rw := newChainRw(ch, log, rollback)
	self := &consensus{rw: rw, rollback: rollback}
	self.mLog = log

	return self
}
