package iavl

import (
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

const (
	IavlStoreAccountKey  = "ANKRCHAINACCOUNT"
	IAvlStoreTransferKey = "ANKRCHAINTRANSFER"
	IAvlStoreMeteringKey = "ANKRCHAINMETERING"
	IAVLStoreAdminKey    = "ANKRCHAINADMIN"

	IavlStoreAccountDefCacheSize  = 10000
	IAvlStoreTransferDefCacheSize = 10000
	IAvlStoreMeteringDefCacheSize = 10000
	IAVLStoreAdminKeyDefCacheSize = 5000

	IAVLStoreAccountKeepVersionNum  = 100
	IAVLStoreTransferKeepVersionNum = 100
	IAVLStoreMeteringKeepVersionNum = 100
	IAVLStoreAdminKeepVersionNum    = 20
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

	dbTran := dbm.NewPrefixDB(db, []byte("ankr:"+IAvlStoreTransferKey+"/"))
	storeMap[IAvlStoreTransferKey] = NewIavlStore(dbTran, IAvlStoreTransferDefCacheSize, IAVLStoreTransferKeepVersionNum, storeLog.With("module", "transferstore"))

	dbMt := dbm.NewPrefixDB(db, []byte("ankr:"+IAvlStoreMeteringKey+"/"))
	storeMap[IAvlStoreMeteringKey] = NewIavlStore(dbMt, IAvlStoreMeteringDefCacheSize, IAVLStoreMeteringKeepVersionNum, storeLog.With("module", "meteringstore"))

	dbAdmin := dbm.NewPrefixDB(db, []byte("ankr:"+IAVLStoreAdminKey+"/"))
	storeMap[IAVLStoreAdminKey] = NewIavlStore(dbAdmin, IAVLStoreAdminKeyDefCacheSize, IAVLStoreAdminKeepVersionNum, storeLog.With("module", "adminstore"))

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




