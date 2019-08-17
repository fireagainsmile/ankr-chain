package iavl

import (
	"fmt"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/merkle"
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

    ProofOpMultiStore = "multistore"
    CommitInfoKey     = "cminfo%d"
    LatestVerKey      = "latestverkey"
)

type IavlStoreMulti struct {
	db           dbm.DB
	storeMap     map[string]*IavlStore
	log          log.Logger
}

type commitInfo struct {
	version int64
	commit  map[string]ankrtypes.CommitID //storeName->CommID
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

func (ms *IavlStoreMulti) IavlStore(storeKey string) *IavlStore {
	if store, ok := ms.storeMap[storeKey]; ok {
		return store
	}

	ms.log.Error("can't find the responding iavlstore", "storeKey", storeKey)

	return nil
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

func (ms *IavlStoreMulti) commitInfo(version int64) *commitInfo {
	infoKey := fmt.Sprintf(CommitInfoKey, version)
	infoV := ms.db.Get([]byte(infoKey))
	if infoV != nil {
		cdc := amino.NewCodec()
		var cmmInfo commitInfo
		cdc.MustUnmarshalBinaryLengthPrefixed(infoV, &cmmInfo)
		return &cmmInfo
	}else {
		ms.log.Error("can't get commit info", "infokey", infoKey)
		return nil
	}
}

func (ms *IavlStoreMulti) setCommitInfo(batch dbm.Batch, version int64, info commitInfo) {
	cdc := amino.NewCodec()
	infoBytes := cdc.MustMarshalBinaryLengthPrefixed(info)
	infoKey := fmt.Sprintf(CommitInfoKey, version)
	batch.Set([]byte(infoKey), infoBytes)
}

func (ms *IavlStoreMulti) latestVersion() int64 {
	verBytes := ms.db.Get([]byte(LatestVerKey))
	if verBytes == nil {
		return 0
	}

	var latestVer int64
	cdc := amino.NewCodec()
	err := cdc.UnmarshalBinaryLengthPrefixed(verBytes, &latestVer)
	if err != nil {
		panic(err)
	}

	return latestVer
}

func (ms *IavlStoreMulti) setLatestVersion(batch dbm.Batch, version int64) {
	cdc := amino.NewCodec()
	latestVerBtest := cdc.MustMarshalBinaryBare(version)
	batch.Set([]byte(LatestVerKey), latestVerBtest)
}

func (ms *IavlStoreMulti) lastCommit() ankrtypes.CommitID {
	latestVer := ms.latestVersion()
    cmmInfos :=  ms.commitInfo(latestVer)
    if cmmInfos != nil {
    	if cmmInfos.version != latestVer {
    		ms.log.Error("error commitinfo, mmInfos.version != latestVer", "latestVer", latestVer, "commitInfoVer", cmmInfos.version)
		}

		hashM := make(map[string][]byte)
		for k, s := range cmmInfos.commit {
			hashM[k] = s.Hash
		}
		reHash := merkle.SimpleHashFromMap(hashM)

		return ankrtypes.CommitID{latestVer, reHash}
	}else {
		ms.log.Error("can't get the latest commitinfo", "latestVer", latestVer)
	}

    return ankrtypes.CommitID{}
}

func (ms *IavlStoreMulti) Commit(version int64) ankrtypes.CommitID {
	var cmmInfo commitInfo

	version += 1

	cmmInfo.version = version
	hashM := make(map[string][]byte)
	for k, s := range ms.storeMap {
		commitID, err := s.Commit()
		if err != nil {
			panic(err)
		}

		hashM[k] = commitID.Hash

		cmmInfo.commit[k] = commitID
	}

	reHash := merkle.SimpleHashFromMap(hashM)

	batch := ms.db.NewBatch()
	ms.setCommitInfo(batch, version, cmmInfo)
	ms.setLatestVersion(batch, version)

	return ankrtypes.CommitID{version, reHash}
}




