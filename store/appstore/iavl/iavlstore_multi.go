package iavl

import (
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

const (
	IavlStoreAccountKey  = "ANKRCHAINACCOUNT"
	IAvlStoreTxKey       = "ANKRCHAINTX"
	IAvlStoreContractKey = "ANKRCHAINCONTRACT"

	IavlStoreAccountDefCacheSize  = 10000
	IAvlStoreTxDefCacheSize       = 10000
	IAvlStoreContractDefCacheSize = 10000


	IAVLStoreAccountKeepVersionNum  = 100
	IAVLStoreTxKeepVersionNum       = 100
	IAVLStoreContractKeepVersionNum = 100
)

type IavlStoreMulti struct {
	db       dbm.DB
	storeMap map[string]*IavlStore
	log      log.Logger
}

func NewIavlStoreMulti(db dbm.DB, storeLog log.Logger) *IavlStoreMulti {

	if db == nil {
		panic("can't create IavlStoreMulti, db nil")
	}

	if storeLog == nil {
		panic("can't create IavlStoreMulti, storeLog nil")
	}

	storeMap := make(map[string]*IavlStore)

	dbAcc := dbm.NewPrefixDB(db, []byte("ankr:"+IavlStoreAccountKey+"/"))
	storeMap[IavlStoreAccountKey] = NewIavlStore(dbAcc, IavlStoreAccountDefCacheSize, IAVLStoreAccountKeepVersionNum, storeLog.With("module", "accountstore"))

	dbTran := dbm.NewPrefixDB(db, []byte("ankr:"+IAvlStoreTxKey+"/"))
	storeMap[IAvlStoreTxKey] = NewIavlStore(dbTran, IAvlStoreTxDefCacheSize, IAVLStoreTxKeepVersionNum, storeLog.With("module", "txstore"))

	dbMt := dbm.NewPrefixDB(db, []byte("ankr:"+IAvlStoreContractKey+"/"))
	storeMap[IAvlStoreContractKey] = NewIavlStore(dbMt, IAvlStoreContractDefCacheSize, IAVLStoreContractKeepVersionNum, storeLog.With("module", "contractstore"))

	return &IavlStoreMulti{db, storeMap, storeLog}
}

// Load the latest versioned tree from disk.
func (ms *IavlStoreMulti) Load() {
	for key, iavlS := range ms.storeMap {
		ver, err := iavlS.Load()
		if err != nil {
			ms.log.Error("Load the latest db failed", "key", key, "err", err)
		}else {
			ms.log.Info("Load the latest db successful", "key", key, "ver", ver)
		}
	}
}




