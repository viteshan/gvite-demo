package consensus

import (
	"time"

	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/consensus/core"
	"github.com/vitelabs/go-vite/log15"
)

var simpleGenesis = time.Unix(1553849738, 0)

var simpleAddrs = genSimpleAddrs()

// vite_360232b0378111b122685a15e612143dc9a89cfa7e803f4b5a hex public key:3fc5224e59433bff4f48c83c0eb4edea0e4c42ea697e04cdec717d03e50d5200
// vite_ce18b99b46c70c8e6bf34177d0c5db956a8c3ea7040a1c1e25 hex public key:e0de77ffdc2719eb1d8e89139da9747bd413bfe59781c43fc078bb37d8cbd77a
func genSimpleAddrs() []types.Address {
	var simpleAddrs []types.Address
	addrs := []string{"vite_360232b0378111b122685a15e612143dc9a89cfa7e803f4b5a",
		"vite_ce18b99b46c70c8e6bf34177d0c5db956a8c3ea7040a1c1e25"}

	for _, v := range addrs {
		addr, err := types.HexToAddress(v)
		if err != nil {
			panic(err)
		}
		simpleAddrs = append(simpleAddrs, addr)
	}
	return simpleAddrs
}

func genSimpleInfo() *core.GroupInfo {
	group := types.ConsensusGroupInfo{
		Gid:                    types.SNAPSHOT_GID,
		NodeCount:              2,
		Interval:               1,
		PerCount:               3,
		RandCount:              1,
		RandRank:               100,
		Repeat:                 1,
		CountingTokenId:        types.CreateTokenTypeId(),
		RegisterConditionId:    0,
		RegisterConditionParam: nil,
		VoteConditionId:        0,
		VoteConditionParam:     nil,
		Owner:                  types.Address{},
		PledgeAmount:           nil,
		WithdrawHeight:         0,
	}

	info := core.NewGroupInfo(simpleGenesis, group)
	return info
}

type simpleCs struct {
	core.GroupInfo
	algo core.Algo

	log log15.Logger
}

func (simple *simpleCs) GetInfo() *core.GroupInfo {
	return &simple.GroupInfo
}

func newSimpleCs(log log15.Logger) *simpleCs {
	cs := &simpleCs{}
	cs.log = log.New("gid", "snapshot")

	cs.GroupInfo = *genSimpleInfo()
	cs.algo = core.NewAlgo(&cs.GroupInfo)
	return cs
}

func (simple *simpleCs) GenProofTime(h uint64) time.Time {
	_, end := simple.Index2Time(h)
	return end
}

func (simple *simpleCs) ElectionTime(t time.Time) (*electionResult, error) {
	index := simple.Time2Index(t)
	return simple.ElectionIndex(index)
}

func (simple *simpleCs) ElectionIndex(index uint64) (*electionResult, error) {
	plans := genElectionResult(&simple.GroupInfo, index, simpleAddrs)
	return plans, nil
}

func (simple *simpleCs) VerifyProducer(address types.Address, t time.Time) (bool, error) {
	electionResult, err := simple.ElectionTime(t)
	if err != nil {
		return false, err
	}

	return simple.verifyProducer(t, address, electionResult), nil
}

func (simple *simpleCs) verifyProducer(t time.Time, address types.Address, result *electionResult) bool {
	if result == nil {
		return false
	}
	for _, plan := range result.Plans {
		if plan.Member == address {
			if plan.STime == t {
				return true
			}
		}
	}
	return false
}
