package onroad

import (
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/ledger"
)

// JudgeGenesis is used for the method ExcludePairTrades.
type JudgeGenesis interface {
	IsGenesisAccountBlock(block types.Hash) bool
}

// ExcludePairTrades is to exclude the trade pairs.
func ExcludePairTrades(chain JudgeGenesis, blockList []*ledger.AccountBlock) map[types.Address][]*ledger.AccountBlock {
	cutMap := make(map[types.Hash]*ledger.AccountBlock)
	for _, block := range blockList {
		if block.IsSendBlock() {
			v, ok := cutMap[block.Hash]
			if ok && v != nil && v.IsReceiveBlock() {
				delete(cutMap, block.Hash)
			} else {
				cutMap[block.Hash] = block
			}
			continue
		}

		if chain.IsGenesisAccountBlock(block.Hash) {
			continue
		}

		// receive block
		v, ok := cutMap[block.FromBlockHash]
		if ok && v != nil && v.IsSendBlock() {
			delete(cutMap, block.FromBlockHash)
		} else {
			cutMap[block.FromBlockHash] = block
		}

		// sendBlockList
		if !types.IsContractAddr(block.AccountAddress) || len(block.SendBlockList) <= 0 {
			continue
		}
		for _, subSend := range block.SendBlockList {
			v, ok := cutMap[subSend.Hash]
			if ok && v != nil && v.IsReceiveBlock() {
				delete(cutMap, subSend.Hash)
			} else {
				cutMap[subSend.Hash] = subSend
			}
		}
	}

	pendingMap := make(map[types.Address][]*ledger.AccountBlock)
	for _, v := range cutMap {
		if v == nil {
			continue
		}
		var addr *types.Address
		if v.IsSendBlock() {
			addr = &v.ToAddress
		} else {
			addr = &v.AccountAddress
		}
		_, ok := pendingMap[*addr]
		if !ok {
			pendingMap[*addr] = make([]*ledger.AccountBlock, 0)
		}
		pendingMap[*addr] = append(pendingMap[*addr], v)
	}
	return pendingMap
}
