// Package vm implements the vite virtual machine
package vm

import (
	"encoding/hex"
	"errors"
	"runtime/debug"

	"github.com/vitelabs/go-vite/common"
	"github.com/vitelabs/go-vite/vm_db"

	"math/big"
	"path/filepath"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/vitelabs/go-vite/log15"

	"github.com/vitelabs/go-vite/common/helper"
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/ledger"
	"github.com/vitelabs/go-vite/monitor"
	"github.com/vitelabs/go-vite/vm/contracts"
	"github.com/vitelabs/go-vite/vm/quota"
	"github.com/vitelabs/go-vite/vm/util"
)

// NodeConfig holds the global status of vm.
type NodeConfig struct {
	isTest         bool
	canTransfer    func(db vm_db.VmDb, tokenTypeId types.TokenTypeId, tokenAmount *big.Int, feeAmount *big.Int) bool
	interpreterLog log15.Logger
	log            log15.Logger
	IsDebug        bool
}

var nodeConfig NodeConfig

// IsTest returns whether node is currently running under a test mode or not.
func IsTest() bool {
	return nodeConfig.isTest
}

// InitVMConfig init global status of vm. It should be
// called when the node started.
func InitVMConfig(isTest bool, isTestParam bool, isDebug bool, datadir string) {
	if isTest {
		nodeConfig = NodeConfig{
			isTest: isTest,
			canTransfer: func(db vm_db.VmDb, tokenTypeId types.TokenTypeId, tokenAmount *big.Int, feeAmount *big.Int) bool {
				return true
			},
		}
	} else {
		nodeConfig = NodeConfig{
			isTest: isTest,
			canTransfer: func(db vm_db.VmDb, tokenTypeId types.TokenTypeId, tokenAmount *big.Int, feeAmount *big.Int) bool {
				if feeAmount.Sign() == 0 {
					b, err := db.GetBalance(&tokenTypeId)
					util.DealWithErr(err)
					return tokenAmount.Cmp(b) <= 0
				}
				if util.IsViteToken(tokenTypeId) {
					balance := new(big.Int).Add(tokenAmount, feeAmount)
					b, err := db.GetBalance(&tokenTypeId)
					util.DealWithErr(err)
					return balance.Cmp(b) <= 0
				}
				amountB, err := db.GetBalance(&tokenTypeId)
				util.DealWithErr(err)
				feeB, err := db.GetBalance(&ledger.ViteTokenId)
				util.DealWithErr(err)
				return tokenAmount.Cmp(amountB) <= 0 && feeAmount.Cmp(feeB) <= 0
			},
		}
	}
	nodeConfig.log = log15.New("module", "vm")
	nodeConfig.interpreterLog = log15.New("module", "vm")
	contracts.InitContractsConfig(isTestParam)
	quota.InitQuotaConfig(isTest, isTestParam)
	nodeConfig.IsDebug = isDebug
	if isDebug {
		initLog(datadir, "dbug")
	}
}

func initLog(dir, lvl string) {
	logLevel, err := log15.LvlFromString(lvl)
	if err != nil {
		logLevel = log15.LvlInfo
	}
	path := filepath.Join(dir, "vmlog", time.Now().Format("2006-01-02T15-04"))
	filename := filepath.Join(path, "vm.log")
	nodeConfig.log.SetHandler(
		log15.LvlFilterHandler(logLevel, log15.StreamHandler(common.MakeDefaultLogger(filename), log15.LogfmtFormat())),
	)
	interpreterFileName := filepath.Join(path, "interpreter.log")
	nodeConfig.interpreterLog.SetHandler(
		log15.LvlFilterHandler(logLevel, log15.StreamHandler(common.MakeDefaultLogger(interpreterFileName), log15.LogfmtFormat())),
	)
}

type vmContext struct {
	sendBlockList []*ledger.AccountBlock
}

// VM holds the runtime information of vite vm and provides
// the necessary tools to run a transfer transaction of a
// call contract transaction. It also provides an offchain
// getter method to read contract storage without a
// transaction.
// The VM instance should never be reused and is not thread
// safe.
type VM struct {
	abort int32
	vmContext
	i            *interpreter
	globalStatus util.GlobalStatus
	reader       util.ConsensusReader
}

// NewVM constructor of VM
func NewVM(cr util.ConsensusReader) *VM {
	return &VM{reader: cr}
}

// GlobalStatus getter
func (vm *VM) GlobalStatus() util.GlobalStatus {
	return vm.globalStatus
}

// ConsensusReader getter
func (vm *VM) ConsensusReader() util.ConsensusReader {
	return vm.reader
}

func printDebugBlockInfo(block *ledger.AccountBlock, result *vm_db.VmAccountBlock, err error) {
	var str string
	if result != nil {
		if result.AccountBlock.IsSendBlock() {
			str = "{SelfAddr: " + result.AccountBlock.AccountAddress.String() + ", " +
				"ToAddr: " + result.AccountBlock.ToAddress.String() + ", " +
				"BlockType: " + strconv.FormatInt(int64(result.AccountBlock.BlockType), 10) + ", " +
				"Quota: " + strconv.FormatUint(result.AccountBlock.Quota, 10) + ", " +
				"Amount: " + result.AccountBlock.Amount.String() + ", " +
				"TokenId: " + result.AccountBlock.TokenId.String() + ", " +
				"Height: " + strconv.FormatUint(result.AccountBlock.Height, 10) + ", " +
				"Data: " + hex.EncodeToString(result.AccountBlock.Data) + ", " +
				"Fee: " + result.AccountBlock.Fee.String() + "}"
		} else {
			if len(result.AccountBlock.SendBlockList) > 0 {
				str = "["
				for _, sendBlock := range result.AccountBlock.SendBlockList {
					str = str + "{ToAddr:" + sendBlock.ToAddress.String() + ", " +
						"BlockType:" + strconv.FormatInt(int64(sendBlock.BlockType), 10) + ", " +
						"Data:" + hex.EncodeToString(sendBlock.Data) + ", " +
						"Amount:" + sendBlock.Amount.String() + ", " +
						"TokenId:" + sendBlock.TokenId.String() + ", " +
						"Fee:" + sendBlock.Fee.String() + "}"
				}
				str = str + "]"
			}
			str = "{SelfAddr: " + result.AccountBlock.AccountAddress.String() + ", " +
				"FromHash: " + result.AccountBlock.FromBlockHash.String() + ", " +
				"BlockType: " + strconv.FormatInt(int64(result.AccountBlock.BlockType), 10) + ", " +
				"Quota: " + strconv.FormatUint(result.AccountBlock.Quota, 10) + ", " +
				"Height: " + strconv.FormatUint(result.AccountBlock.Height, 10) + ", " +
				"Data: " + hex.EncodeToString(result.AccountBlock.Data) + ", " +
				"SendBlockList: " + str + "}"
		}
	}
	nodeConfig.log.Info("vm run stop",
		"blockType", block.BlockType,
		"address", block.AccountAddress.String(),
		"height", block.Height,
		"fromHash", block.FromBlockHash.String(),
		"err", err,
		"block", str,
	)
}

func getContractMeta(db vm_db.VmDb) *ledger.ContractMeta {
	ok, err := db.IsContractAccount()
	util.DealWithErr(err)
	if !ok {
		return nil
	}
	meta, err := db.GetContractMeta()
	util.DealWithErr(err)
	return meta
}

// RunV2 method executes an account block, performs balance change and storage change, returns execution result
func (vm *VM) RunV2(db vm_db.VmDb, block *ledger.AccountBlock, sendBlock *ledger.AccountBlock, status util.GlobalStatus) (vmAccountBlock *vm_db.VmAccountBlock, isRetry bool, err error) {
	defer monitor.LogTimerConsuming([]string{"vm", "run"}, time.Now())
	defer func() {
		db.Finish()
		if nodeConfig.IsDebug {
			printDebugBlockInfo(block, vmAccountBlock, err)
		}
	}()
	if nodeConfig.IsDebug {
		nodeConfig.log.Info("vm run start",
			"blockType", block.BlockType,
			"address", block.AccountAddress.String(),
			"height", block.Height, ""+
				"fromHash", block.FromBlockHash.String())
	}
	blockcopy := block.Copy()
	sb, err := db.LatestSnapshotBlock()
	util.DealWithErr(err)
	vm.i = newInterpreter(sb.Height, false)
	vm.globalStatus = status
	switch block.BlockType {
	case ledger.BlockTypeReceive, ledger.BlockTypeReceiveError:
		blockcopy.Data = nil
		contractMeta := getContractMeta(db)
		if sendBlock.BlockType == ledger.BlockTypeSendCreate {
			return vm.receiveCreate(db, blockcopy, sendBlock, quota.CalcCreateQuota(sendBlock.Fee), contractMeta)
		} else if sendBlock.BlockType == ledger.BlockTypeSendCall || sendBlock.BlockType == ledger.BlockTypeSendReward {
			return vm.receiveCall(db, blockcopy, sendBlock, contractMeta)
		} else if sendBlock.BlockType == ledger.BlockTypeSendRefund {
			return vm.receiveRefund(db, blockcopy, sendBlock, contractMeta)
		}
	case ledger.BlockTypeSendCreate:
		quotaTotal, quotaAddition, err := quota.CalcQuotaForBlock(
			db,
			block.AccountAddress,
			getPledgeAmount(db),
			block.Difficulty)
		if err != nil {
			return nil, noRetry, err
		}
		vmAccountBlock, err = vm.sendCreate(db, blockcopy, true, quotaTotal, quotaAddition)
		if err != nil {
			return nil, noRetry, err
		}
		return vmAccountBlock, noRetry, nil
	case ledger.BlockTypeSendCall:
		quotaTotal, quotaAddition, err := quota.CalcQuotaForBlock(
			db,
			block.AccountAddress,
			getPledgeAmount(db),
			block.Difficulty)
		if err != nil {
			return nil, noRetry, err
		}
		vmAccountBlock, err = vm.sendCall(db, blockcopy, true, quotaTotal, quotaAddition)
		if err != nil {
			return nil, noRetry, err
		}
		return vmAccountBlock, noRetry, nil
	case ledger.BlockTypeSendReward, ledger.BlockTypeSendRefund:
		return nil, noRetry, util.ErrTransactionTypeNotSupport
	}
	return nil, noRetry, util.ErrTransactionTypeNotSupport
}

// Cancel method stops the running contract receive
func (vm *VM) Cancel() {
	atomic.StoreInt32(&vm.abort, 1)
}

// send contract create transaction, create address, sub balance and service fee
func (vm *VM) sendCreate(db vm_db.VmDb, block *ledger.AccountBlock, useQuota bool, quotaTotal, quotaAddition uint64) (*vm_db.VmAccountBlock, error) {
	defer monitor.LogTimerConsuming([]string{"vm", "sendCreate"}, time.Now())
	// check can make transaction
	quotaLeft := quotaTotal
	if useQuota {
		cost, err := gasNormalSendCall(block)
		if err != nil {
			return nil, err
		}
		quotaLeft, err = util.UseQuota(quotaLeft, cost)
		if err != nil {
			return nil, err
		}
	}
	var err error
	block.Fee, err = calcContractFee(block.Data)
	if err != nil {
		return nil, err
	}

	gid := util.GetGidFromCreateContractData(block.Data)
	if gid == types.SNAPSHOT_GID {
		return nil, util.ErrInvalidMethodParam
	}

	contractType := util.GetContractTypeFromCreateContractData(block.Data)
	if !util.IsExistContractType(contractType) {
		return nil, util.ErrInvalidMethodParam
	}

	confirmTime := util.GetConfirmTimeFromCreateContractData(block.Data)
	if confirmTime < confirmTimeMin || confirmTime > confirmTimeMax {
		return nil, util.ErrInvalidConfirmTime
	}

	quotaRatio := util.GetQuotaRatioFromCreateContractData(block.Data)
	if !util.IsValidQuotaRatio(quotaRatio) {
		return nil, util.ErrInvalidQuotaRatio
	}

	if ContainsStatusCode(util.GetCodeFromCreateContractData(block.Data)) && confirmTime <= 0 {
		return nil, util.ErrInvalidConfirmTime
	}

	if !nodeConfig.canTransfer(db, block.TokenId, block.Amount, block.Fee) {
		return nil, util.ErrInsufficientBalance
	}

	contractAddr := util.NewContractAddress(
		block.AccountAddress,
		block.Height,
		block.PrevHash)

	block.ToAddress = contractAddr
	// sub balance and service fee
	util.SubBalance(db, &block.TokenId, block.Amount)
	util.SubBalance(db, &ledger.ViteTokenId, block.Fee)
	q, qUsed := util.CalcQuotaUsed(useQuota, quotaTotal, quotaAddition, quotaLeft, nil)
	vm.updateBlock(db, block, nil, q, qUsed)
	db.SetContractMeta(contractAddr, &ledger.ContractMeta{Gid: gid, SendConfirmedTimes: confirmTime, QuotaRatio: quotaRatio})
	return &vm_db.VmAccountBlock{block, db}, nil
}

// receive contract create transaction, create contract account, run initialization code, set contract code, do send blocks
func (vm *VM) receiveCreate(db vm_db.VmDb, block *ledger.AccountBlock, sendBlock *ledger.AccountBlock, quotaTotal uint64, meta *ledger.ContractMeta) (*vm_db.VmAccountBlock, bool, error) {
	defer monitor.LogTimerConsuming([]string{"vm", "receiveCreate"}, time.Now())

	quotaLeft := quotaTotal
	prev, err := db.PrevAccountBlock()
	util.DealWithErr(err)
	if prev != nil {
		return nil, noRetry, util.ErrAddressCollision
	}
	// check can make transaction
	cost, err := gasReceiveCreate(block, meta)
	if err != nil {
		return nil, noRetry, err
	}
	quotaLeft, err = util.UseQuota(quotaLeft, cost)
	if err != nil {
		return nil, noRetry, err
	}

	// create contract account and add balance
	util.AddBalance(db, &sendBlock.TokenId, sendBlock.Amount)

	// init contract state_bak and set contract code
	initCode := util.GetCodeFromCreateContractData(sendBlock.Data)
	c := newContract(block, db, sendBlock, initCode, quotaLeft)
	c.setCallCode(block.AccountAddress, initCode)
	code, err := c.run(vm)
	if err == nil && len(code) <= maxCodeSize {
		code := util.PackContractCode(util.GetContractTypeFromCreateContractData(sendBlock.Data), code)
		codeCost := uint64(len(code)) * contractCodeGas
		c.quotaLeft, err = util.UseQuota(c.quotaLeft, codeCost)
		if err == nil {
			db.SetContractCode(code)
			vm.updateBlock(db, block, err, 0, 0)
			db, err = vm.doSendBlockList(db)
			if err == nil {
				block.Data = getReceiveCallData(db, err, 0)
				return mergeReceiveBlock(db, block, vm.sendBlockList), noRetry, nil
			}
		}
	}
	vm.revert(db)

	// try refund
	vm.updateBlock(db, block, err, 0, 0)
	if sendBlock.Amount.Sign() > 0 {
		vm.vmContext.AppendBlock(
			util.MakeSendBlock(
				block.AccountAddress,
				sendBlock.AccountAddress,
				ledger.BlockTypeSendRefund,
				new(big.Int).Set(sendBlock.Amount),
				sendBlock.TokenId,
				[]byte{}))
		util.AddBalance(db, &sendBlock.TokenId, sendBlock.Amount)
		var refundErr error
		if db, refundErr = vm.doSendBlockList(db); refundErr == nil {
			block.Data = getReceiveCallData(db, err, 0)
			return mergeReceiveBlock(db, block, vm.sendBlockList), noRetry, err
		}
		monitor.LogEvent("vm", "impossibleReceiveError")
		nodeConfig.log.Error("Impossible receive error", "err", refundErr, "fromhash", sendBlock.Hash)
		return nil, retry, err
	}
	block.Data = getReceiveCallData(db, err, 0)
	return &vm_db.VmAccountBlock{block, db}, noRetry, err
}

func mergeReceiveBlock(db vm_db.VmDb, receiveBlock *ledger.AccountBlock, sendBlockList []*ledger.AccountBlock) *vm_db.VmAccountBlock {
	receiveBlock.SendBlockList = sendBlockList
	return &vm_db.VmAccountBlock{receiveBlock, db}
}

func (vm *VM) sendCall(db vm_db.VmDb, block *ledger.AccountBlock, useQuota bool, quotaTotal, quotaAddition uint64) (*vm_db.VmAccountBlock, error) {
	defer monitor.LogTimerConsuming([]string{"vm", "sendCall"}, time.Now())
	// check can make transaction
	quotaLeft := quotaTotal
	if p, ok, err := contracts.GetBuiltinContractMethod(block.ToAddress, block.Data); ok {
		if err != nil {
			return nil, err
		}
		block.Fee, err = p.GetFee(block)
		if err != nil {
			return nil, err
		}
		if !nodeConfig.canTransfer(db, block.TokenId, block.Amount, block.Fee) {
			return nil, util.ErrInsufficientBalance
		}
		if useQuota {
			cost, err := p.GetSendQuota(block.Data)
			if err != nil {
				return nil, err
			}
			quotaLeft, err = util.UseQuota(quotaLeft, cost)
			if err != nil {
				return nil, err
			}
		}
		err = p.DoSend(db, block)
		if err != nil {
			return nil, err
		}
		util.SubBalance(db, &block.TokenId, block.Amount)
		util.SubBalance(db, &ledger.ViteTokenId, block.Fee)
	} else {
		block.Fee = helper.Big0
		if useQuota {
			quotaLeft, err = useQuotaForSend(block, db, quotaLeft)
			if err != nil {
				return nil, err
			}
		}
		if !nodeConfig.canTransfer(db, block.TokenId, block.Amount, block.Fee) {
			return nil, util.ErrInsufficientBalance
		}
		util.SubBalance(db, &block.TokenId, block.Amount)
	}
	q, qUsed := util.CalcQuotaUsed(useQuota, quotaTotal, quotaAddition, quotaLeft, nil)
	vm.updateBlock(db, block, nil, q, qUsed)
	return &vm_db.VmAccountBlock{block, db}, nil
}

var (
	resultSuccess  = byte(0)
	resultFail     = byte(1)
	resultDepthErr = byte(2)
)

func getReceiveCallData(db vm_db.VmDb, err error, qUsed uint64) []byte {
	if err == nil {
		return append(db.GetReceiptHash().Bytes(), resultSuccess)
	} else if err == util.ErrDepth {
		return append(db.GetReceiptHash().Bytes(), resultDepthErr)
	} else {
		return append(db.GetReceiptHash().Bytes(), resultFail)
	}
}

func (vm *VM) receiveCall(db vm_db.VmDb, block *ledger.AccountBlock, sendBlock *ledger.AccountBlock, meta *ledger.ContractMeta) (*vm_db.VmAccountBlock, bool, error) {
	defer monitor.LogTimerConsuming([]string{"vm", "receiveCall"}, time.Now())

	if checkDepth(db, sendBlock) {
		util.AddBalance(db, &sendBlock.TokenId, sendBlock.Amount)
		block.Data = getReceiveCallData(db, util.ErrDepth, 0)
		vm.updateBlock(db, block, util.ErrDepth, 0, 0)
		return &vm_db.VmAccountBlock{block, db}, noRetry, util.ErrDepth
	}
	if p, ok, _ := contracts.GetBuiltinContractMethod(block.AccountAddress, sendBlock.Data); ok {
		util.AddBalance(db, &sendBlock.TokenId, sendBlock.Amount)
		blockListToSend, err := p.DoReceive(db, block, sendBlock, vm)
		if err == nil {
			vm.updateBlock(db, block, err, 0, 0)
			vm.vmContext.sendBlockList = blockListToSend
			if db, err = vm.doSendBlockList(db); err == nil {
				block.Data = getReceiveCallData(db, err, 0)
				return mergeReceiveBlock(db, block, vm.sendBlockList), noRetry, nil
			}
		}
		vm.revert(db)
		refundFlag := false
		refundData, needRefund := p.GetRefundData()
		refundFlag = doRefund(vm, db, block, sendBlock, refundData, needRefund, ledger.BlockTypeSendCall)
		vm.updateBlock(db, block, err, 0, 0)
		if refundFlag {
			var refundErr error
			if db, refundErr = vm.doSendBlockList(db); refundErr != nil {
				monitor.LogEvent("vm", "impossibleReceiveError")
				nodeConfig.log.Error("Impossible receive error", "err", refundErr, "fromhash", sendBlock.Hash)
				return nil, retry, err
			}
			block.Data = getReceiveCallData(db, err, 0)
			return mergeReceiveBlock(db, block, vm.sendBlockList), noRetry, err
		}
		block.Data = getReceiveCallData(db, err, 0)
		return &vm_db.VmAccountBlock{block, db}, noRetry, err
	}
	// check can make transaction
	quotaTotal, quotaAddition, err := quota.CalcQuotaForBlock(
		db,
		block.AccountAddress,
		getPledgeAmount(db),
		block.Difficulty)
	util.DealWithErr(err)
	quotaLeft := quotaTotal
	cost, err := gasReceive(block, meta)
	if err != nil {
		return nil, noRetry, err
	}
	quotaLeft, err = util.UseQuota(quotaLeft, cost)
	if err != nil {
		return nil, retry, err
	}
	// add balance, create account if not exist
	util.AddBalance(db, &sendBlock.TokenId, sendBlock.Amount)
	// do transfer transaction if account code size is zero
	if meta == nil {
		q, qUsed := util.CalcQuotaUsed(true, quotaTotal, quotaAddition, quotaLeft, nil)
		vm.updateBlock(db, block, nil, q, qUsed)
		return &vm_db.VmAccountBlock{block, db}, noRetry, nil
	}
	// run code
	_, code := util.GetContractCode(db, &block.AccountAddress, nil)
	c := newContract(block, db, sendBlock, sendBlock.Data, quotaLeft)
	c.setCallCode(block.AccountAddress, code)
	_, err = c.run(vm)
	if err == nil {
		q, qUsed := util.CalcQuotaUsed(true, quotaTotal, quotaAddition, c.quotaLeft, nil)
		vm.updateBlock(db, block, err, q, qUsed)
		db, err = vm.doSendBlockList(db)
		if err == nil {
			block.Data = getReceiveCallData(db, err, qUsed)
			return mergeReceiveBlock(db, block, vm.sendBlockList), noRetry, nil
		}
	}

	if err == util.ErrNoReliableStatus {
		return nil, retry, err
	}

	vm.revert(db)

	if err == util.ErrOutOfQuota {
		unConfirmedList := db.GetUnconfirmedBlocks(*db.Address())
		if len(unConfirmedList) > 0 {
			// Contract receive out of quota, current block is not first unconfirmed block, retry next snapshotBlock
			return nil, retry, err
		}
		// Contract receive out of quota, current block is first unconfirmed block, refund with no quota
		refundFlag := doRefund(vm, db, block, sendBlock, []byte{}, false, ledger.BlockTypeSendRefund)
		q, qUsed := util.CalcQuotaUsed(true, quotaTotal, quotaAddition, c.quotaLeft, err)
		vm.updateBlock(db, block, err, q, qUsed)
		if refundFlag {
			var refundErr error
			if db, refundErr = vm.doSendBlockList(db); refundErr != nil {
				monitor.LogEvent("vm", "impossibleReceiveError")
				nodeConfig.log.Error("Impossible receive error", "err", refundErr, "fromhash", sendBlock.Hash)
				return nil, retry, err
			}
			block.Data = getReceiveCallData(db, err, qUsed)
			return mergeReceiveBlock(db, block, vm.sendBlockList), noRetry, err
		}
		block.Data = getReceiveCallData(db, err, qUsed)
		return &vm_db.VmAccountBlock{block, db}, noRetry, err
	}

	refundFlag := doRefund(vm, db, block, sendBlock, []byte{}, false, ledger.BlockTypeSendRefund)
	q, qUsed := util.CalcQuotaUsed(true, quotaTotal, quotaAddition, c.quotaLeft, err)
	vm.updateBlock(db, block, err, q, qUsed)
	if refundFlag {
		var refundErr error
		if db, refundErr = vm.doSendBlockList(db); refundErr != nil {
			monitor.LogEvent("vm", "impossibleReceiveError")
			nodeConfig.log.Error("Impossible receive error", "err", refundErr, "fromhash", sendBlock.Hash)
			return nil, retry, err
		}
		block.Data = getReceiveCallData(db, err, qUsed)
		return mergeReceiveBlock(db, block, vm.sendBlockList), noRetry, err
	}
	block.Data = getReceiveCallData(db, err, qUsed)
	return &vm_db.VmAccountBlock{block, db}, noRetry, err
}

func doRefund(vm *VM, db vm_db.VmDb, block *ledger.AccountBlock, sendBlock *ledger.AccountBlock, refundData []byte, needRefund bool, refundBlockType byte) bool {
	refundFlag := false
	if sendBlock.Amount.Sign() > 0 && sendBlock.Fee.Sign() > 0 && sendBlock.TokenId == ledger.ViteTokenId {
		refundAmount := new(big.Int).Add(sendBlock.Amount, sendBlock.Fee)
		vm.vmContext.AppendBlock(
			util.MakeSendBlock(
				block.AccountAddress,
				sendBlock.AccountAddress,
				refundBlockType,
				refundAmount,
				ledger.ViteTokenId,
				refundData))
		util.AddBalance(db, &ledger.ViteTokenId, refundAmount)
		refundFlag = true
	} else {
		if sendBlock.Amount.Sign() > 0 {
			vm.vmContext.AppendBlock(
				util.MakeSendBlock(
					block.AccountAddress,
					sendBlock.AccountAddress,
					refundBlockType,
					new(big.Int).Set(sendBlock.Amount),
					sendBlock.TokenId,
					refundData))
			util.AddBalance(db, &sendBlock.TokenId, sendBlock.Amount)
			refundFlag = true
		}
		if sendBlock.Fee.Sign() > 0 {
			vm.vmContext.AppendBlock(
				util.MakeSendBlock(
					block.AccountAddress,
					sendBlock.AccountAddress,
					refundBlockType,
					new(big.Int).Set(sendBlock.Fee),
					ledger.ViteTokenId,
					refundData))
			util.AddBalance(db, &ledger.ViteTokenId, sendBlock.Fee)
			refundFlag = true
		}
	}
	if !refundFlag && needRefund {
		vm.vmContext.AppendBlock(
			util.MakeSendBlock(
				block.AccountAddress,
				sendBlock.AccountAddress,
				ledger.BlockTypeSendCall,
				big.NewInt(0),
				ledger.ViteTokenId,
				refundData))
		refundFlag = true
	}
	return refundFlag
}

func (vm *VM) sendReward(db vm_db.VmDb, block *ledger.AccountBlock, useQuota bool, quotaTotal, quotaAddition uint64) (*vm_db.VmAccountBlock, error) {
	defer monitor.LogTimerConsuming([]string{"vm", "sendReward"}, time.Now())

	// check can make transaction
	quotaLeft := quotaTotal
	if useQuota {
		var err error
		quotaLeft, err = useQuotaForSend(block, db, quotaLeft)
		if err != nil {
			return nil, err
		}
	}
	if block.AccountAddress != types.AddressConsensusGroup &&
		block.AccountAddress != types.AddressMintage {
		return nil, util.ErrInvalidMethodParam
	}
	q, qUsed := util.CalcQuotaUsed(useQuota, quotaTotal, quotaAddition, quotaLeft, nil)
	vm.updateBlock(db, block, nil, q, qUsed)
	return &vm_db.VmAccountBlock{block, db}, nil
}

func (vm *VM) sendRefund(db vm_db.VmDb, block *ledger.AccountBlock, useQuota bool, quotaTotal, quotaAddition uint64) (*vm_db.VmAccountBlock, error) {
	defer monitor.LogTimerConsuming([]string{"vm", "sendRefund"}, time.Now())
	block.Fee = helper.Big0
	quotaLeft := quotaTotal
	if useQuota {
		var err error
		quotaLeft, err = useQuotaForSend(block, db, quotaLeft)
		if err != nil {
			return nil, err
		}
	}
	if !nodeConfig.canTransfer(db, block.TokenId, block.Amount, block.Fee) {
		return nil, util.ErrInsufficientBalance
	}
	util.SubBalance(db, &block.TokenId, block.Amount)
	q, qUsed := util.CalcQuotaUsed(useQuota, quotaTotal, quotaAddition, quotaLeft, nil)
	vm.updateBlock(db, block, nil, q, qUsed)
	return &vm_db.VmAccountBlock{block, db}, nil
}

func (vm *VM) receiveRefund(db vm_db.VmDb, block *ledger.AccountBlock, sendBlock *ledger.AccountBlock, meta *ledger.ContractMeta) (*vm_db.VmAccountBlock, bool, error) {
	defer monitor.LogTimerConsuming([]string{"vm", "receiveRefund"}, time.Now())
	// check can make transaction
	quotaTotal, quotaAddition, err := quota.CalcQuotaForBlock(
		db,
		block.AccountAddress,
		getPledgeAmount(db),
		block.Difficulty)
	util.DealWithErr(err)
	quotaLeft := quotaTotal
	cost, err := gasReceive(block, meta)
	if err != nil {
		return nil, noRetry, err
	}
	quotaLeft, err = util.UseQuota(quotaLeft, cost)
	if err != nil {
		return nil, retry, err
	}
	util.AddBalance(db, &sendBlock.TokenId, sendBlock.Amount)
	q, qUsed := util.CalcQuotaUsed(true, quotaTotal, quotaAddition, quotaLeft, nil)
	vm.updateBlock(db, block, nil, q, qUsed)
	return &vm_db.VmAccountBlock{block, db}, noRetry, nil
}

func (vm *VM) delegateCall(contractAddr types.Address, data []byte, c *contract) (ret []byte, err error) {
	_, code := util.GetContractCode(c.db, &contractAddr, vm.globalStatus)
	if len(code) > 0 {
		cNew := newContract(c.block, c.db, c.sendBlock, c.data, c.quotaLeft)
		cNew.setCallCode(contractAddr, code)
		ret, err = cNew.run(vm)
		c.quotaLeft = cNew.quotaLeft
		return ret, err
	}
	return nil, nil
}

func (vm *VM) updateBlock(db vm_db.VmDb, block *ledger.AccountBlock, err error, q, qUsed uint64) {
	block.Quota = q
	block.QuotaUsed = qUsed
	if block.IsReceiveBlock() {
		block.LogHash = db.GetLogListHash()
		if err == util.ErrOutOfQuota {
			block.BlockType = ledger.BlockTypeReceiveError
		} else {
			block.BlockType = ledger.BlockTypeReceive
		}
	}
}

func (vm *VM) doSendBlockList(db vm_db.VmDb) (newDb vm_db.VmDb, err error) {
	if len(vm.sendBlockList) == 0 {
		return db, nil
	}
	for i, block := range vm.sendBlockList {
		var sendBlock *vm_db.VmAccountBlock
		switch block.BlockType {
		case ledger.BlockTypeSendCall:
			sendBlock, err = vm.sendCall(db, block, false, 0, 0)
			if err != nil {
				return db, err
			}
		case ledger.BlockTypeSendReward:
			sendBlock, err = vm.sendReward(db, block, false, 0, 0)
			if err != nil {
				return db, err
			}
		case ledger.BlockTypeSendRefund:
			sendBlock, err = vm.sendRefund(db, block, false, 0, 0)
			if err != nil {
				return db, err
			}
		}
		vm.sendBlockList[i] = sendBlock.AccountBlock
		db = sendBlock.VmDb
	}
	return db, nil
}

func (vm *VM) revert(db vm_db.VmDb) {
	vm.sendBlockList = nil
	db.Reset()
}

// AppendBlock method append a send block to send block list of a contract receive block
func (context *vmContext) AppendBlock(block *ledger.AccountBlock) {
	context.sendBlockList = append(context.sendBlockList, block)
}

func calcContractFee(data []byte) (*big.Int, error) {
	return createContractFee, nil
}

func checkDepth(db vm_db.VmDb, sendBlock *ledger.AccountBlock) bool {
	depth, err := db.GetCallDepth(&sendBlock.Hash)
	util.DealWithErr(err)
	return depth >= callDepth
}

// OffChainReader read contract storage without tx
func (vm *VM) OffChainReader(db vm_db.VmDb, code []byte, data []byte) (result []byte, err error) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			nodeConfig.log.Error("offchain reader panic",
				"err", err,
				"addr", db.Address(),
				"code", hex.EncodeToString(code),
				"data", hex.EncodeToString(data))
			result = nil
			err = errors.New("offchain reader panic")
		}
	}()
	sb, err := db.LatestSnapshotBlock()
	if err != nil {
		return nil, err
	}
	vm.i = newInterpreter(sb.Height, true)
	c := newContract(&ledger.AccountBlock{AccountAddress: *db.Address()}, db, &ledger.AccountBlock{ToAddress: *db.Address()}, data, offChainReaderGas)
	c.setCallCode(*db.Address(), code)
	return c.run(vm)
}

func getPledgeAmount(db vm_db.VmDb) *big.Int {
	pledgeAmount, err := db.GetPledgeBeneficialAmount(db.Address())
	util.DealWithErr(err)
	return pledgeAmount
}

func useQuotaForSend(block *ledger.AccountBlock, db vm_db.VmDb, quotaLeft uint64) (uint64, error) {
	cost, err := gasNormalSendCall(block)
	if err != nil {
		return quotaLeft, err
	}
	quotaRatio, err := getQuotaRatioForS(db, block.ToAddress)
	if err != nil {
		return quotaLeft, err
	}
	cost, err = util.MultipleCost(cost, quotaRatio)
	if err != nil {
		return quotaLeft, err
	}
	quotaLeft, err = util.UseQuota(quotaLeft, cost)
	return quotaLeft, err
}
