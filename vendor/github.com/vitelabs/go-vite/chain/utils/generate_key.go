package chain_utils

import (
	"encoding/binary"
	"github.com/vitelabs/go-vite/common/types"
)

// ====== index db ======
func CreateAccountBlockHashKey(blockHash *types.Hash) []byte {
	key := make([]byte, 0, 1+types.HashSize)
	key = append(key, AccountBlockHashKeyPrefix)
	key = append(key, blockHash.Bytes()...)
	return key
}

func CreateAccountBlockHeightKey(addr *types.Address, height uint64) []byte {
	key := make([]byte, 0, 1+types.AddressSize+8)

	key = append(key, AccountBlockHeightKeyPrefix)
	key = append(key, addr.Bytes()...)
	key = append(key, Uint64ToBytes(height)...)
	return key
}

func CreateReceiveKey(sendBlockHash *types.Hash) []byte {
	key := make([]byte, 0, 1+types.HashSize)
	key = append(key, ReceiveKeyPrefix)
	key = append(key, sendBlockHash.Bytes()...)
	return key
}

func CreateConfirmHeightKey(addr *types.Address, height uint64) []byte {
	key := make([]byte, 0, 1+types.AddressSize+8)
	key = append(key, ConfirmHeightKeyPrefix)
	key = append(key, addr.Bytes()...)
	key = append(key, Uint64ToBytes(height)...)
	return key
}

func CreateAccountAddressKey(addr *types.Address) []byte {
	addrBytes := addr.Bytes()
	key := make([]byte, 0, 1+types.AddressSize)
	key = append(key, AccountAddressKeyPrefix)
	key = append(key, addrBytes...)
	return key
}

func CreateOnRoadKey(toAddr types.Address, blockHash types.Hash) []byte {
	key := make([]byte, 0, 1+types.AddressSize+types.HashSize)
	key = append(key, OnRoadKeyPrefix)
	key = append(key, toAddr.Bytes()...)
	key = append(key, blockHash.Bytes()...)

	return key
}

func CreateAccountIdKey(accountId uint64) []byte {
	key := make([]byte, 0, 9)
	key = append(key, AccountIdKeyPrefix)
	key = append(key, Uint64ToBytes(accountId)...)

	return key
}

func CreateSnapshotBlockHashKey(snapshotBlockHash *types.Hash) []byte {
	key := make([]byte, 0, 1+types.HashSize)
	key = append(key, SnapshotBlockHashKeyPrefix)
	key = append(key, snapshotBlockHash.Bytes()...)
	return key
}

func CreateSnapshotBlockHeightKey(snapshotBlockHeight uint64) []byte {
	key := make([]byte, 0, 9)
	key = append(key, SnapshotBlockHeightKeyPrefix)
	key = append(key, Uint64ToBytes(snapshotBlockHeight)...)
	return key
}

// ====== state db ======

func CreateStorageValueKeyPrefix(address *types.Address, prefix []byte) []byte {
	keySize := 1 + types.AddressSize + len(prefix)
	key := make([]byte, 0, keySize)

	key = append(key, StorageKeyPrefix)
	key = append(key, address.Bytes()...)
	key = append(key, prefix...)
	return key
}

func CreateStorageValueKey(address *types.Address, storageKey []byte) []byte {
	keySize := 1 + types.AddressSize + types.HashSize + 1
	key := make([]byte, keySize)
	key[0] = StorageKeyPrefix

	copy(key[1:1+types.AddressSize], address.Bytes())
	copy(key[1+types.AddressSize:1+types.AddressSize+types.HashSize], storageKey)
	key[keySize-1] = byte(len(storageKey))

	return key
}

func CreateHistoryStorageValueKeyPrefix(address *types.Address, prefix []byte) []byte {
	keySize := 1 + types.AddressSize + len(prefix)
	key := make([]byte, 0, keySize)

	key = append(key, StorageHistoryKeyPrefix)
	key = append(key, address.Bytes()...)
	key = append(key, prefix...)
	return key
}

func CreateHistoryStorageValueKey(address *types.Address, storageKey []byte, snapshotHeight uint64) []byte {
	keySize := 1 + types.AddressSize + types.HashSize + 1 + 8
	key := make([]byte, keySize)
	key[0] = StorageHistoryKeyPrefix

	copy(key[1:1+types.AddressSize], address.Bytes())
	copy(key[1+types.AddressSize:1+types.AddressSize+types.HashSize], storageKey)

	key[keySize-9] = byte(len(storageKey))
	binary.BigEndian.PutUint64(key[keySize-8:], snapshotHeight)

	return key
}

func CreateBalanceKeyPrefix(address types.Address) []byte {
	key := make([]byte, 1+types.AddressSize)
	key[0] = BalanceKeyPrefix
	copy(key[1:types.AddressSize+1], address.Bytes())

	return key
}
func CreateBalanceKey(address types.Address, tokenTypeId types.TokenTypeId) []byte {
	key := make([]byte, 1+types.AddressSize+types.TokenTypeIdSize)
	key[0] = BalanceKeyPrefix

	copy(key[1:types.AddressSize+1], address.Bytes())
	copy(key[types.AddressSize+1:], tokenTypeId.Bytes())

	return key
}

func CreateHistoryBalanceKey(address types.Address, tokenTypeId types.TokenTypeId, snapshotHeight uint64) []byte {
	key := make([]byte, 1+types.AddressSize+types.TokenTypeIdSize+8)
	key[0] = BalanceHistoryKeyPrefix

	copy(key[1:types.AddressSize+1], address.Bytes())
	copy(key[types.AddressSize+1:], tokenTypeId.Bytes())
	binary.BigEndian.PutUint64(key[len(key)-8:], snapshotHeight)

	return key
}

func CreateCodeKey(address types.Address) []byte {
	keySize := 1 + types.AddressSize

	key := make([]byte, keySize)

	key[0] = CodeKeyPrefix

	copy(key[1:], address.Bytes())

	return key
}

func CreateContractMetaKey(address types.Address) []byte {
	keySize := 1 + types.AddressSize

	key := make([]byte, keySize)

	key[0] = ContractMetaKeyPrefix

	copy(key[1:], address.Bytes())

	return key
}
func CreateGidContractKey(gid types.Gid, address *types.Address) []byte {
	key := make([]byte, 0, 1+types.GidSize+types.AddressSize)

	key = append(key, GidContractKeyPrefix)

	key = append(key, gid.Bytes()...)
	key = append(key, address.Bytes()...)

	return key
}

func CreateGidContractPrefixKey(gid *types.Gid) []byte {
	key := make([]byte, 0, 1+types.GidSize)

	key = append(key, GidContractKeyPrefix)
	key = append(key, gid.Bytes()...)

	return key
}

func CreateVmLogListKey(logHash *types.Hash) []byte {
	key := make([]byte, 1+types.HashSize)

	key[0] = VmLogListKeyPrefix

	copy(key[1:], logHash.Bytes())

	return key
}

func CreateCallDepthKey(blockHash types.Hash) []byte {
	key := make([]byte, 0, 1+types.HashSize)
	key = append(key, CallDepthKeyPrefix)
	key = append(key, blockHash.Bytes()...)
	return key
}

// ====== state redo ======

func CreateRedoSnapshot(snapshotHeight uint64) []byte {
	key := make([]byte, 0, 1+8)
	key = append(key, SnapshotKeyPrefix)
	key = append(key, Uint64ToBytes(snapshotHeight)...)
	return key
}
