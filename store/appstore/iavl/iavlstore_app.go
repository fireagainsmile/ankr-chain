package iavl

import (
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/merkle"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

const (
	AccountKey = "acckey"
	CertKey    = "certkey"
)

type IavlStoreApp struct {
	iavlSM          *IavlStoreMulti
	lastCommitID    ankrtypes.CommitID
	storeLog        log.Logger
	accStoreLocker  sync.RWMutex
	certStoreLocker sync.RWMutex
}

func NewIavlStoreApp(dbDir string, storeLog log.Logger) *IavlStoreApp {
	kvPath := filepath.Join(dbDir, "kvstore.db")
	isKVPathExist, err := common.PathExists(kvPath)
	if err != nil {
		panic(err)
	}

	var kvDB dbm.DB
	var lcmmID ankrtypes.CommitID
	if isKVPathExist {
		kvDB, err = dbm.NewGoLevelDB("kvstore.db", dbDir)
		if err != nil {
			panic(err)
		}

		oldState := appstore.LoadState(kvDB)
		lcmmID.Version = oldState.Height
		lcmmID.Hash    = oldState.AppHash
	}

	db, err := dbm.NewGoLevelDB("appstore.db", dbDir)
	if err != nil {
		panic(err)
	}

	iavlSM := NewIavlStoreMulti(db, storeLog)

	if !isKVPathExist {
		lcmmID = iavlSM.lastCommit()
	}

	iavlSApp := &IavlStoreApp{iavlSM: iavlSM, lastCommitID: lcmmID, storeLog: storeLog}

	iavlSM.storeMap[IavlStoreAccountKey].Set([]byte(AccountKey), []byte(""))
	iavlSM.storeMap[IAvlStoreTxKey].Set([]byte(CertKey), []byte(""))

	if isKVPathExist {
		iavlSApp.Prefixed(kvDB, kvPath)
	}

	return iavlSApp
}

func (sp* IavlStoreApp) Prefixed(kvDB dbm.DB, kvPath string) error {
	var iavlStore *IavlStore
	it := kvDB.Iterator(nil, nil)
	defer it.Close()
	if it != nil && it.Valid(){
		it.Next()
		for it.Valid() {
			if len(it.Key()) >= len(ankrtypes.AccountBlancePrefix) && string(it.Key()[0:len(ankrtypes.AccountBlancePrefix)]) == ankrtypes.AccountBlancePrefix {
				iavlStore = sp.iavlSM.IavlStore(IavlStoreAccountKey)
				sp.SetBalance(it.Key(), it.Value())
			}else {
				iavlStore = sp.iavlSM.IavlStore(IAvlStoreTxKey)
				if len(it.Key()) >= len(ankrtypes.CertPrefix) && string(it.Key()[0:len(ankrtypes.CertPrefix)]) == ankrtypes.CertPrefix {
					sp.SetCertKey(it.Key(), it.Value())
				} else {
					iavlStore.Set(it.Key(), it.Value())
				}
			}
			it.Next()
		}
	}

	os.RemoveAll(kvPath)

	return nil
}

func (sp* IavlStoreApp) updateAccount(addr string) {
	sp.accStoreLocker.Lock()
	defer sp.accStoreLocker.Unlock()

	accs, err := sp.iavlSM.IavlStore(IavlStoreAccountKey).Get([]byte(AccountKey))
	if err == nil {
		accs = append(accs, []byte(";" + addr)...)
		sp.iavlSM.IavlStore(IavlStoreAccountKey).Set([]byte(AccountKey), accs)
	}else {
		sp.storeLog.Error("can't get the AccountKey value", "err",  err)
	}
}

func (sp *IavlStoreApp) SetBalance(key []byte, val []byte) {
	if !sp.iavlSM.IavlStore(IavlStoreAccountKey).Has(key) {
        accP := string(key)
        address := strings.Split(accP, ":")
		sp.updateAccount(address[1])
	}

	sp.iavlSM.IavlStore(IavlStoreAccountKey).Set(key, val)
}

func (sp *IavlStoreApp) Balance(key []byte) []byte {
	balV, err := sp.iavlSM.IavlStore(IavlStoreAccountKey).Get(key)
	if err != nil {
		sp.storeLog.Error("can't get balance", "key", key, "err", err)
		balV = nil
	}

	return balV
}

func (sp *IavlStoreApp) AccountList() []byte {
	accs, err := sp.iavlSM.IavlStore(IavlStoreAccountKey).Get([]byte(AccountKey))
	if err == nil {
		return accs
	}

	return nil
}

func (sp *IavlStoreApp) LastCommit() *ankrtypes.CommitID{
	return &sp.lastCommitID
}

func (sp *IavlStoreApp) Commit() types.ResponseCommit {
    commitID := sp.iavlSM.Commit(sp.lastCommitID.Version)

	sp.lastCommitID.Hash = sp.lastCommitID.Hash[0:0]

    sp.lastCommitID.Version = commitID.Version
	sp.lastCommitID.Hash    = append(sp.lastCommitID.Hash, commitID.Hash...)

	return types.ResponseCommit{Data: commitID.Hash}
}

func (sp *IavlStoreApp) parsePath(path string)(storeName string, subPath string) {
	if path == "" || !strings.HasPrefix(path, "/"){
		sp.storeLog.Error("invalid path", "path", path)
		return "", ""
	}

	pathSegs := strings.Split(path[1:], "/")
	storeName = pathSegs[0]
	if len(pathSegs) == 2 {
		subPath = "/" + pathSegs[1]
	}

	return
}

func (sp *IavlStoreApp) Query(reqQuery types.RequestQuery) (resQuery types.ResponseQuery) {
	var value []byte
	resQuery.Log = "exists"

	storeName, _ := sp.parsePath(reqQuery.Path)
	if reqQuery.Path != "" {
		storeName = IAvlStoreTxKey
	}

	if reqQuery.Prove {
		qVer := reqQuery.Height
		if qVer == 0 {
			qVer = sp.lastCommitID.Version
		}
		commInfo := sp.iavlSM.commitInfo(qVer)

		if commInfo != nil {
			cdc := amino.NewCodec()
			infoBytes := cdc.MustMarshalBinaryLengthPrefixed(commInfo)
			pOP := merkle.ProofOp{
				Type: ProofOpMultiStore,
				Key:  []byte(storeName),
				Data: infoBytes,
			}
			resQuery.Proof.Ops = append(resQuery.Proof.Ops, pOP)
		}
	}

	resQuery.Key = reqQuery.Data
	if string(reqQuery.Data[:3]) == ankrtypes.AccountBlancePrefix[:3] {
		value = sp.Balance(reqQuery.Data)
		trxGetBalanceSlices := strings.Split(string(value), ":")
		if len(trxGetBalanceSlices) == 1 {
			_, err := new(big.Int).SetString(string(value), 10)
			if !err {
				resQuery.Log = "internal error, value format incorrect, single value"
			}
		}else if len(trxGetBalanceSlices) == 2 {
			_, berr := new(big.Int).SetString(trxGetBalanceSlices[0], 10)
			if !berr {
				resQuery.Log = "internal error, value format incorrect, first value"
			} else {
				_, err := strconv.ParseInt(string(trxGetBalanceSlices[1]), 10, 64)
				if err != nil {
					resQuery.Log = "internal error, value format incorrect, second value"
				}
			}
		} else {
			resQuery.Log = "internal error, value format incorrect, extra value"
		}
	} else if len(reqQuery.Data) >= len(ankrtypes.AllAccountsPrefix) && string(reqQuery.Data[:len(ankrtypes.AllAccountsPrefix)]) == ankrtypes.AllAccountsPrefix {
		value = sp.AccountList()
	} else if len(reqQuery.Data) >= len(ankrtypes.AllCrtsPrefix) && string(reqQuery.Data[:len(ankrtypes.AllCrtsPrefix)]) == ankrtypes.AllCrtsPrefix {
		value = sp.CertKeyList()
	} else {
		value, _ = sp.iavlSM.IavlStore(IAvlStoreTxKey).Get(reqQuery.Data)
	}

	resQuery.Value = value

	if value == nil {
		resQuery.Log = "does not exist"
	}

	return
}

func (sp *IavlStoreApp) updateCertKey(dcS string){
	sp.certStoreLocker.Lock()
	defer sp.certStoreLocker.Unlock()

	certs, err := sp.iavlSM.IavlStore(IAvlStoreTxKey).Get([]byte(CertKey))
	if err == nil {
		certs = append(certs, []byte(";" + dcS)...)
		sp.iavlSM.IavlStore(IAvlStoreTxKey).Set([]byte(CertKey), certs)
	}else {
		sp.storeLog.Error("can't get the CertKey value", "err",  err)
	}
}

func (sp *IavlStoreApp) SetCertKey(key []byte, val []byte) {
	if !sp.iavlSM.IavlStore(IAvlStoreTxKey).Has(key) {
		dcS := strings.Split(string(key), ":")[1]
		sp.updateCertKey(dcS)
	}

	sp.iavlSM.IavlStore(IAvlStoreTxKey).Set(key, val)
}

func (sp *IavlStoreApp) CertKey(key []byte) []byte {
	valBytes, err :=  sp.iavlSM.IavlStore(IAvlStoreTxKey).Get(key)
	if err != nil {
		sp.storeLog.Error("can't get the key's value", "key", key)
		valBytes = nil
	}

	return valBytes
}

func (sp *IavlStoreApp) DeleteCertKey(key []byte) {
	sp.iavlSM.IavlStore(IAvlStoreTxKey).Remove(key)
}

func (sp *IavlStoreApp) CertKeyList() []byte {
	certs, err := sp.iavlSM.IavlStore(IAvlStoreTxKey).Get([]byte(CertKey))
	if err == nil {
		return certs
	}

	return nil
}

func (sp *IavlStoreApp) Get(key []byte) []byte {
	val, err := sp.iavlSM.IavlStore(IAvlStoreTxKey).Get(key)
	if err != nil {
		sp.storeLog.Error("can't get the key value", "key", string(key))
		val = nil
	}

	return val
}

func (sp *IavlStoreApp) Set(key []byte, val []byte) {
	sp.iavlSM.IavlStore(IAvlStoreTxKey).Set(key, val)
}

func (sp *IavlStoreApp) Delete(key []byte) {
	sp.iavlSM.IavlStore(IAvlStoreTxKey).Remove(key)
}

func (sp *IavlStoreApp) Has(key []byte) bool {
	return sp.iavlSM.IavlStore(IAvlStoreTxKey).Has(key)
}

func (sp *IavlStoreApp) Height() int64 {
	return sp.lastCommitID.Version
}

func (sp *IavlStoreApp) APPHash() []byte {
	return sp.lastCommitID.Hash
}

func (sp *IavlStoreApp) DB() dbm.DB {
	return sp.iavlSM.db
}


