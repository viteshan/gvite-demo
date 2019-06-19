package api

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/vitelabs/go-vite/chain"
	"github.com/vitelabs/go-vite/common/math"
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/ledger"
	"github.com/vitelabs/go-vite/onroad"
	"github.com/vitelabs/go-vite/vite"
	"strconv"
)

var MaxBatchQuery = 10

type PublicOnroadApi struct {
	api *PrivateOnroadApi
}

func NewPublicOnroadApi(vite *vite.Vite) *PublicOnroadApi {
	return &PublicOnroadApi{
		api: NewPrivateOnroadApi(vite),
	}
}

func (pub PublicOnroadApi) String() string {
	return "PublicOnroadApi"
}

func (pub PublicOnroadApi) GetOnroadBlocksByAddress(address types.Address, index, count uint64) ([]*AccountBlock, error) {
	if count > math.MaxUint16+1 {
		return nil, errors.New(fmt.Sprintf("maximum number per page allowed is %d", math.MaxUint16+1))
	}
	return pub.api.GetOnroadBlocksByAddress(address, index, count)
}

func (pub PublicOnroadApi) GetOnroadInfoByAddress(address types.Address) (*RpcAccountInfo, error) {
	return pub.api.GetOnroadInfoByAddress(address)
}

type PrivateOnroadApi struct {
	manager *onroad.Manager
}

func NewPrivateOnroadApi(vite *vite.Vite) *PrivateOnroadApi {
	return &PrivateOnroadApi{
		manager: vite.OnRoad(),
	}
}

func (pri PrivateOnroadApi) String() string {
	return "PrivateOnroadApi"
}

func (pri PrivateOnroadApi) GetContractAddrListByGid(gid types.Gid) ([]types.Address, error) {
	return pri.manager.Chain().GetContractList(gid)
}

func (pri PrivateOnroadApi) GetOnroadBlocksByAddress(address types.Address, index, count uint64) ([]*AccountBlock, error) {
	log.Info("GetOnroadBlocksByAddress", "addr", address, "index", index, "count", count)
	blockList, err := pri.manager.Chain().GetOnRoadBlocksByAddr(address, int(index), int(count))
	if err != nil {
		return nil, err
	}
	a := make([]*AccountBlock, len(blockList))
	sum := 0
	for k, v := range blockList {
		if v != nil {
			accountBlock, e := ledgerToRpcBlock(pri.manager.Chain(), v)
			if e != nil {
				return nil, e
			}
			a[k] = accountBlock
			sum++
		}
	}
	return a[:sum], nil
}

func (pri PrivateOnroadApi) GetOnroadInfoByAddress(address types.Address) (*RpcAccountInfo, error) {
	log.Info("GetAccountOnroadInfo", "addr", address)
	info, e := pri.manager.Chain().GetAccountOnRoadInfo(address)
	if e != nil || info == nil {
		return nil, e
	}
	r := onroadInfoToRpcAccountInfo(pri.manager.Chain(), info)
	return r, nil

}

type OnroadPagingQuery struct {
	Addr types.Address `json:"addr"`

	PageNum   uint64 `json:"pageNum"`
	PageCount uint64 `json:"pageCount"`
}

func (pri PrivateOnroadApi) GetOnroadBlocksInBatch(queryList []OnroadPagingQuery) (map[types.Address][]*AccountBlock, error) {
	if len(queryList) <= 0 {
		return nil, nil
	}
	if len(queryList) > MaxBatchQuery {
		return nil, errors.New(fmt.Sprintf("the maximum number of queries allowed is %d", MaxBatchQuery))
	}
	resultMap := make(map[types.Address][]*AccountBlock)
	for _, q := range queryList {
		if l, ok := resultMap[q.Addr]; ok && l != nil {
			continue
		}
		if q.PageCount > math.MaxUint8+1 {
			return nil, errors.New(fmt.Sprintf("maximum number per page allowed is %d", math.MaxUint8+1))
		}
		blockList, err := pri.GetOnroadBlocksByAddress(q.Addr, q.PageNum, q.PageCount)
		if err != nil {
			return nil, err
		}
		if len(blockList) <= 0 {
			continue
		}
		resultMap[q.Addr] = blockList
	}
	return resultMap, nil
}

func (pri PrivateOnroadApi) GetOnroadInfoInBatch(addrList []types.Address) ([]*RpcAccountInfo, error) {
	if len(addrList) <= 0 {
		return nil, nil
	}
	if len(addrList) > MaxBatchQuery {
		return nil, errors.New(fmt.Sprintf("the maximum number of queries allowed is %d", MaxBatchQuery))
	}
	// Remove duplicate
	addrs := make([]types.Address, 0)
	for _, v1 := range addrList {
		var isExist = false
		for _, v2 := range addrs {
			if v2 == v1 {
				isExist = true
				break
			}
		}
		if !isExist {
			addrs = append(addrs, v1)
		}
	}

	resultList := make([]*RpcAccountInfo, 0)
	for _, addr := range addrs {
		info, err := pri.GetOnroadInfoByAddress(addr)
		if err != nil {
			return nil, err
		}
		if info == nil {
			continue
		}
		resultList = append(resultList, info)
	}
	return resultList, nil
}

func onroadInfoToRpcAccountInfo(chain chain.Chain, info *ledger.AccountInfo) *RpcAccountInfo {
	var r RpcAccountInfo
	r.AccountAddress = info.AccountAddress
	r.TotalNumber = strconv.FormatUint(info.TotalNumber, 10)
	r.TokenBalanceInfoMap = make(map[types.TokenTypeId]*RpcTokenBalanceInfo)

	for tti, v := range info.TokenBalanceInfoMap {
		if v != nil {
			number := strconv.FormatUint(v.Number, 10)
			tinfo, _ := chain.GetTokenInfoById(tti)
			b := &RpcTokenBalanceInfo{
				TokenInfo:   RawTokenInfoToRpc(tinfo, tti),
				TotalAmount: v.TotalAmount.String(),
				Number:      &number,
			}
			r.TokenBalanceInfoMap[tti] = b
		}
	}
	return &r
}

func (pri PrivateOnroadApi) GetContractOnRoadTotalNum(addr types.Address, gid *types.Gid) (uint64, error) {
	if !types.IsContractAddr(addr) {
		return 0, errors.New("Address must be the type of Contract.")
	}

	var g types.Gid
	if gid == nil {
		g = types.DELEGATE_GID
	} else {
		g = *gid
	}

	num, err := pri.manager.GetOnRoadTotalNumByAddr(g, addr)
	if err != nil {
		return 0, err
	}
	log.Info("GetContractOnRoadTotalNum", "gid", gid, "addr", addr, "num", num)
	return num, nil
}

func (pri PrivateOnroadApi) GetContractOnRoadFrontBlocks(addr types.Address, gid *types.Gid) ([]*AccountBlock, error) {
	if !types.IsContractAddr(addr) {
		return nil, errors.New("Address must be the type of Contract.")
	}
	var g types.Gid
	if gid == nil {
		g = types.DELEGATE_GID
	} else {
		g = *gid
	}

	blockList, err := pri.manager.GetAllCallersFrontOnRoad(g, addr)
	if err != nil {
		return nil, err
	}
	log.Info("GetContractOnRoadFrontBlocks", "gid", gid, "addr", addr, "len", len(blockList))
	rpcBlockList := make([]*AccountBlock, len(blockList))
	sum := 0
	for k, v := range blockList {
		if v != nil {
			accountBlock, e := ledgerToRpcBlock(pri.manager.Chain(), v)
			if e != nil {
				return nil, e
			}
			rpcBlockList[k] = accountBlock
			sum++
		}
	}
	return rpcBlockList[:sum], nil
}
