package ledger

import (
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"github.com/vitelabs/go-vite/common"
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/crypto"
	"github.com/vitelabs/go-vite/crypto/ed25519"
	"github.com/vitelabs/go-vite/log15"
	"github.com/vitelabs/go-vite/vitepb"
	"math/big"
)

var accountBlockLog = log15.New("module", "ledger/account_block")

const (
	BlockTypeSendCreate   = byte(1) // send
	BlockTypeSendCall     = byte(2) // send
	BlockTypeSendReward   = byte(3) // send
	BlockTypeReceive      = byte(4) // receive
	BlockTypeReceiveError = byte(5) // receive
	BlockTypeSendRefund   = byte(6) // send

	BlockTypeGenesisReceive = byte(7) // receive
)

type AccountBlock struct {
	BlockType byte       `json:"blockType"` // 1
	Hash      types.Hash `json:"hash"`
	PrevHash  types.Hash `json:"prevHash"` // 2
	Height    uint64     `json:"height"`   // 3

	AccountAddress types.Address `json:"accountAddress"` // 4

	producer  *types.Address    // for cache
	PublicKey ed25519.PublicKey `json:"publicKey"`
	ToAddress types.Address     `json:"toAddress"` // 5

	Amount *big.Int `json:"amount"` // 6	padding 32 bytes

	TokenId types.TokenTypeId `json:"tokenId"` // 7

	FromBlockHash types.Hash `json:"fromBlockHash"` // 8

	Data []byte `json:"data"` // 9	hash

	Quota uint64 `json:"quota"` // quotaUsed = quota + pow quota

	QuotaUsed uint64 `json:"quotaUsed"` // quotaUsed = quota + pow quota

	Fee *big.Int `json:"fee"` // 10 padding 32 bytes

	LogHash *types.Hash `json:"logHash"` // 11

	Difficulty *big.Int `json:"difficulty"`

	Nonce []byte `json:"nonce"` // 12 padding 8 bytes

	SendBlockList []*AccountBlock `json:"sendBlockList"` // 13

	Signature []byte `json:"signature"`
}

func (ab *AccountBlock) Copy() *AccountBlock {
	newAb := *ab

	if ab.Amount != nil {
		newAb.Amount = new(big.Int).Set(ab.Amount)
	}

	if ab.Fee != nil {
		newAb.Fee = new(big.Int).Set(ab.Fee)
	}

	newAb.Data = make([]byte, len(ab.Data))
	copy(newAb.Data, ab.Data)

	if ab.LogHash != nil {
		logHash := *ab.LogHash
		newAb.LogHash = &logHash
	}

	if ab.Difficulty != nil {
		newAb.Difficulty = new(big.Int).Set(ab.Difficulty)
	}

	if len(ab.Nonce) > 0 {
		newAb.Nonce = make([]byte, len(ab.Nonce))
		copy(newAb.Nonce, ab.Nonce)
	}

	if len(ab.Signature) > 0 {
		newAb.Signature = make([]byte, len(ab.Signature))
		copy(newAb.Signature, ab.Signature)
	}

	for _, sendBlock := range ab.SendBlockList {
		newAb.SendBlockList = append(newAb.SendBlockList, sendBlock.Copy())
	}
	return &newAb
}

func (ab *AccountBlock) hashSourceLength() int {
	// 1, 2, 3 , 4
	size := 1 + types.HashSize + 8 + types.AddressSize

	if ab.IsSendBlock() {
		// 5, 6, 7
		size += types.AddressSize + len(ab.Amount.Bytes()) + types.TokenTypeIdSize
	} else {
		// 8
		size += types.HashSize
	}

	// 9
	if len(ab.Data) > 0 {
		size += types.HashSize
	}

	// 10
	size += types.HashSize

	// 11
	if ab.LogHash != nil {
		size += types.HashSize
	}

	// 12, 13
	size += 8 + types.HashSize*len(ab.SendBlockList)

	return size
}

func (ab *AccountBlock) hashSource(extraByte []byte) []byte {
	source := make([]byte, 0, ab.hashSourceLength()+len(extraByte))
	// BlockType
	source = append(source, ab.BlockType)

	// PrevHash
	source = append(source, ab.PrevHash.Bytes()...)

	// Height
	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, ab.Height)
	source = append(source, heightBytes...)

	// AccountAddress
	source = append(source, ab.AccountAddress.Bytes()...)

	if ab.IsSendBlock() {
		// ToAddress
		source = append(source, ab.ToAddress.Bytes()...)
		// Amount(fixed 32 bytes, left padding)
		source = append(source, common.LeftPadBytes(ab.Amount.Bytes(), 32)...)
		// TokenId
		source = append(source, ab.TokenId.Bytes()...)
	} else {
		// FromBlockHash
		source = append(source, ab.FromBlockHash.Bytes()...)
	}

	// Data
	if len(ab.Data) > 0 {
		dataHashBytes := crypto.Hash256(ab.Data)
		source = append(source, dataHashBytes...)
	}

	// Fee(fixed 32 bytes, left padding)
	var feeBytes []byte
	if ab.Fee != nil {
		feeBytes = ab.Fee.Bytes()
	}
	source = append(source, common.LeftPadBytes(feeBytes, 32)...)

	// LogHash
	if ab.LogHash != nil {
		source = append(source, ab.LogHash.Bytes()...)
	}

	// Nonce(fixed 8 bytes, left padding)
	source = append(source, common.LeftPadBytes(ab.Nonce, 8)...)

	// Send block hash list
	for _, sendBlock := range ab.SendBlockList {
		source = append(source, sendBlock.Hash.Bytes()...)
	}

	source = append(source, extraByte...)
	return source
}

func (ab *AccountBlock) ComputeHash() types.Hash {
	source := ab.hashSource(nil)

	hash, _ := types.BytesToHash(crypto.Hash256(source))

	return hash
}

func (ab *AccountBlock) ComputeSendHash(hostBlock *AccountBlock, index uint8) types.Hash {

	extraBytes := make([]byte, 0, types.HashSize+8+1)
	// prev hash
	extraBytes = append(extraBytes, hostBlock.PrevHash.Bytes()...)

	// height
	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, hostBlock.Height)
	extraBytes = append(extraBytes, heightBytes...)

	// index
	extraBytes = append(extraBytes, index)

	source := ab.hashSource(extraBytes)

	hash, _ := types.BytesToHash(crypto.Hash256(source))

	return hash
}

func (ab *AccountBlock) Producer() types.Address {
	if ab.producer == nil {
		producer := types.PubkeyToAddress(ab.PublicKey)
		ab.producer = &producer
	}

	return *ab.producer
}

func (ab *AccountBlock) VerifySignature() bool {
	isVerified, verifyErr := crypto.VerifySig(ab.PublicKey, ab.Hash.Bytes(), ab.Signature)
	if verifyErr != nil {
		accountBlockLog.Error("crypto.VerifySig failed, error is "+verifyErr.Error(), "method", "VerifySignature")
	}
	return isVerified
}

func (ab *AccountBlock) Proto() *vitepb.AccountBlock {
	pb := &vitepb.AccountBlock{}
	// 1
	pb.BlockType = vitepb.AccountBlock_BlockType(ab.BlockType)
	// 2
	pb.Hash = ab.Hash.Bytes()
	// 3
	pb.Height = ab.Height
	// 4
	if ab.Height > 1 {
		pb.PrevHash = ab.PrevHash.Bytes()
	}
	// 5
	pb.AccountAddress = ab.AccountAddress.Bytes()
	// 6
	pb.PublicKey = ab.PublicKey
	if ab.IsSendBlock() {
		// 7
		pb.ToAddress = ab.ToAddress.Bytes()
	} else {
		// 10
		pb.FromBlockHash = ab.FromBlockHash.Bytes()
	}

	if ab.IsSendBlock() || ab.BlockType == BlockTypeGenesisReceive {
		// 8
		pb.Amount = ab.Amount.Bytes()
		// 9
		pb.TokenId = ab.TokenId.Bytes()
	}

	// 11
	pb.Data = ab.Data

	// 12
	pb.Quota = ab.Quota

	// 20
	pb.QuotaUsed = ab.QuotaUsed

	if ab.Fee != nil {
		// 13
		pb.Fee = ab.Fee.Bytes()
	}

	if ab.LogHash != nil {
		// 14
		pb.LogHash = ab.LogHash.Bytes()
	}

	if ab.Difficulty != nil {
		// 15
		pb.Difficulty = ab.Difficulty.Bytes()
	}
	// 16
	pb.Nonce = ab.Nonce
	// 17
	pb.SendBlockList = make([]*vitepb.AccountBlock, 0, len(ab.SendBlockList))
	for _, sendBlock := range ab.SendBlockList {
		pb.SendBlockList = append(pb.SendBlockList, sendBlock.Proto())
	}
	// 18
	pb.Signature = ab.Signature
	return pb
}

func (ab *AccountBlock) DeProto(pb *vitepb.AccountBlock) error {
	var err error
	// 1
	ab.BlockType = byte(pb.BlockType)
	// 2
	ab.Hash, _ = types.BytesToHash(pb.Hash)
	// 3
	ab.Height = pb.Height
	// 4
	if ab.Height > 1 {
		ab.PrevHash, _ = types.BytesToHash(pb.PrevHash)
	}
	// 5
	if ab.AccountAddress, err = types.BytesToAddress(pb.AccountAddress); err != nil {
		return err
	}
	// 6
	ab.PublicKey = pb.PublicKey

	if ab.IsSendBlock() {
		// 7
		if ab.ToAddress, err = types.BytesToAddress(pb.ToAddress); err != nil {
			return err
		}

	} else {
		// 10
		if ab.FromBlockHash, err = types.BytesToHash(pb.FromBlockHash); err != nil {
			return err
		}
	}

	if ab.IsSendBlock() || ab.BlockType == BlockTypeGenesisReceive {

		// 9
		if ab.TokenId, err = types.BytesToTokenTypeId(pb.TokenId); err != nil {
			return err
		}
		// 8
		ab.Amount = big.NewInt(0)
		if len(pb.Amount) > 0 {
			ab.Amount.SetBytes(pb.Amount)
		}
	}

	// 11
	ab.Data = pb.Data

	// 12
	ab.Quota = pb.Quota

	// 20
	ab.QuotaUsed = pb.QuotaUsed

	// 13
	ab.Fee = big.NewInt(0)
	if len(pb.Fee) > 0 {
		ab.Fee.SetBytes(pb.Fee)
	}

	// 14
	if len(pb.LogHash) > 0 {
		logHash, err := types.BytesToHash(pb.LogHash)
		if err != nil {
			return err
		}

		ab.LogHash = &logHash
	}

	// 15
	if len(pb.Difficulty) > 0 {
		ab.Difficulty = new(big.Int).SetBytes(pb.Difficulty)
	}
	// 16
	ab.Nonce = pb.Nonce

	// 17
	ab.SendBlockList = make([]*AccountBlock, 0, len(pb.SendBlockList))
	for _, pbSendBlock := range pb.SendBlockList {
		sendBlock := &AccountBlock{}
		if err := sendBlock.DeProto(pbSendBlock); err != nil {
			return err
		}
		ab.SendBlockList = append(ab.SendBlockList, sendBlock)
	}
	// 18
	ab.Signature = pb.Signature
	return nil
}

func (ab *AccountBlock) Serialize() ([]byte, error) {
	return proto.Marshal(ab.Proto())
}

func (ab *AccountBlock) Deserialize(buf []byte) error {
	pb := &vitepb.AccountBlock{}
	if err := proto.Unmarshal(buf, pb); err != nil {
		return err
	}
	ab.DeProto(pb)
	return nil
}

func (ab *AccountBlock) IsSendBlock() bool {
	return IsSendBlock(ab.BlockType)
}

func IsSendBlock(blockType byte) bool {
	return blockType == BlockTypeSendCreate ||
		blockType == BlockTypeSendCall ||
		blockType == BlockTypeSendReward ||
		blockType == BlockTypeSendRefund
}

func (ab *AccountBlock) IsReceiveBlock() bool {
	return IsReceiveBlock(ab.BlockType)
}

func IsReceiveBlock(blockType byte) bool {
	return blockType == BlockTypeReceive ||
		blockType == BlockTypeReceiveError ||
		blockType == BlockTypeGenesisReceive
}
