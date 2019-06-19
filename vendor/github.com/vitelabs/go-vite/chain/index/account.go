package chain_index

import (
	"encoding/binary"
	"github.com/vitelabs/go-vite/chain/utils"
	"github.com/vitelabs/go-vite/common/db/xleveldb/util"
	"github.com/vitelabs/go-vite/common/helper"
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/interfaces"
	"sync/atomic"
)

func (iDB *IndexDB) HasAccount(addr types.Address) (result bool, returnErr error) {
	_, ok := iDB.accountCache.Get(addr)
	if ok {
		return ok, nil
	} else {
		defer func() {
			if result {
				iDB.accountCache.Add(addr, nil)
			}
		}()
	}

	return iDB.store.Has(chain_utils.CreateAccountAddressKey(&addr))
}

func (iDB *IndexDB) GetAccountId(addr *types.Address) (uint64, error) {
	key := chain_utils.CreateAccountAddressKey(addr)
	value, err := iDB.store.Get(key)
	if err != nil {
		return 0, err
	}

	if len(value) <= 0 {
		return 0, nil
	}
	return chain_utils.BytesToUint64(value), nil
}

func (iDB *IndexDB) GetAccountAddress(accountId uint64) (*types.Address, error) {
	key := chain_utils.CreateAccountIdKey(accountId)
	value, err := iDB.store.Get(key)

	if err != nil {
		return nil, err
	}
	if len(value) <= 0 {
		return nil, nil
	}

	addr, err := types.BytesToAddress(value)
	if err != nil {
		return nil, err
	}
	return &addr, nil
}

func (iDB *IndexDB) IterateAccounts(iterateFunc func(addr types.Address, accountId uint64, err error) bool) {
	iter := iDB.store.NewIterator(util.BytesPrefix([]byte{chain_utils.AccountIdKeyPrefix}))
	defer iter.Release()

	for iter.Next() {
		accountId := binary.BigEndian.Uint64(iter.Key()[1:])
		account, err := types.BytesToAddress(iter.Value())

		if !iterateFunc(account, accountId, err) {
			break
		}
		if err != nil {
			break
		}
	}
}

func (iDB *IndexDB) createAccount(batch interfaces.Batch, addr *types.Address) uint64 {
	newAccountId := atomic.AddUint64(&iDB.latestAccountId, 1)

	batch.Put(chain_utils.CreateAccountAddressKey(addr), chain_utils.Uint64ToBytes(newAccountId))
	batch.Put(chain_utils.CreateAccountIdKey(newAccountId), addr.Bytes())
	return newAccountId
}

func (iDB *IndexDB) queryLatestAccountId() (uint64, error) {
	startKey := chain_utils.CreateAccountIdKey(1)
	endKey := chain_utils.CreateAccountIdKey(helper.MaxUint64)

	iter := iDB.store.NewIterator(&util.Range{Start: startKey, Limit: endKey})
	defer iter.Release()

	var latestAccountId uint64
	if iter.Last() {
		latestAccountId = chain_utils.BytesToUint64(iter.Key()[1:])
	}
	if err := iter.Error(); err != nil {
		return 0, err
	}

	return latestAccountId, nil
}
