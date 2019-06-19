package util

import (
	"bytes"
	"encoding/hex"
	"github.com/vitelabs/go-vite/common/helper"
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/ledger"
	"github.com/vitelabs/go-vite/vm/abi"
	"math/big"
	"sort"
)

var (
	AttovPerVite = big.NewInt(1e18)
)

func IsViteToken(tokenId types.TokenTypeId) bool {
	return tokenId == ledger.ViteTokenId
}
func IsSnapshotGid(gid types.Gid) bool {
	return gid == types.SNAPSHOT_GID
}
func IsDelegateGid(gid types.Gid) bool {
	return gid == types.DELEGATE_GID
}

func MakeSendBlock(fromAddress types.Address, toAddress types.Address, blockType byte, amount *big.Int, tokenId types.TokenTypeId, data []byte) *ledger.AccountBlock {
	return &ledger.AccountBlock{
		AccountAddress: fromAddress,
		ToAddress:      toAddress,
		BlockType:      blockType,
		Amount:         amount,
		TokenId:        tokenId,
		Data:           data,
		Fee:            big.NewInt(0),
	}
}

var (
	SolidityPPContractType = []byte{1}
	contractTypeSize       = 1
	confirmTimeSize        = 1
	quotaRatioSize         = 1
)

func IsValidQuotaRatio(quotaRatio uint8) bool {
	return quotaRatio >= 10 && quotaRatio <= 100
}

func GetCreateContractData(bytecode []byte, contractType []byte, confirmTimes uint8, quotaRatio uint8, gid types.Gid) []byte {
	return helper.JoinBytes(gid.Bytes(), contractType, []byte{confirmTimes}, []byte{quotaRatio}, bytecode)
}

func GetGidFromCreateContractData(data []byte) types.Gid {
	gid, _ := types.BytesToGid(data[:types.GidSize])
	return gid
}

func GetCodeFromCreateContractData(data []byte) []byte {
	return data[types.GidSize+contractTypeSize+confirmTimeSize+quotaRatioSize:]
}
func GetContractTypeFromCreateContractData(data []byte) []byte {
	return data[types.GidSize : types.GidSize+contractTypeSize]
}
func IsExistContractType(contractType []byte) bool {
	if bytes.Equal(contractType, SolidityPPContractType) {
		return true
	}
	return false
}
func GetConfirmTimeFromCreateContractData(data []byte) uint8 {
	return uint8(data[types.GidSize+contractTypeSize])
}
func GetQuotaRatioFromCreateContractData(data []byte) uint8 {
	return uint8(data[types.GidSize+contractTypeSize+confirmTimeSize])
}

func PackContractCode(contractType, code []byte) []byte {
	return helper.JoinBytes(contractType, code)
}

type CommonDb interface {
	Address() *types.Address
	IsContractAccount() (bool, error)
	GetContractCode() ([]byte, error)
	GetContractCodeBySnapshotBlock(addr *types.Address, snapshotBlock *ledger.SnapshotBlock) ([]byte, error)
}

func GetContractCode(db CommonDb, addr *types.Address, status GlobalStatus) ([]byte, []byte) {
	var code []byte
	var err error
	if *db.Address() == *addr {
		code, err = db.GetContractCode()
	} else {
		code, err = db.GetContractCodeBySnapshotBlock(addr, status.SnapshotBlock())
	}
	DealWithErr(err)
	if len(code) > 0 {
		return code[:contractTypeSize], code[contractTypeSize:]
	}
	return nil, nil
}

func NewContractAddress(accountAddress types.Address, accountBlockHeight uint64, prevBlockHash types.Hash) types.Address {
	return types.CreateContractAddress(
		accountAddress.Bytes(),
		helper.LeftPadBytes(new(big.Int).SetUint64(accountBlockHeight).Bytes(), 8),
		prevBlockHash.Bytes())
}

func PrintMap(m map[string][]byte) string {
	var result string
	if len(m) > 0 {
		var keys []string
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			result += hex.EncodeToString([]byte(k)) + "=>" + hex.EncodeToString(m[k]) + ", "
		}
		result = result[:len(result)-2]
	}
	return result
}

func IsUserAccount(addr types.Address) bool {
	return !types.IsContractAddr(addr)
}

func NewLog(c abi.ABIContract, name string, params ...interface{}) *ledger.VmLog {
	topics, data, _ := c.PackEvent(name, params...)
	return &ledger.VmLog{Topics: topics, Data: data}
}
