package iavl

import (
	"errors"
	"fmt"
	"github.com/Ankr-network/ankr-chain/account"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/Ankr-network/ankr-chain/common"
	apscomm "github.com/Ankr-network/ankr-chain/store/appstore/common"
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
	cdc             *amino.Codec
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
		kvDB, err = dbm.NewGoLevelDB("kvstore", dbDir)
		if err != nil {
			panic(err)
		}

		oldState := apscomm.LoadState(kvDB)
		lcmmID.Version = oldState.Height
		lcmmID.Hash    = oldState.AppHash
	}

	db, err := dbm.NewGoLevelDB("appstore", dbDir)
	if err != nil {
		panic(err)
	}

	iavlSM := NewIavlStoreMulti(db, storeLog)

	if !isKVPathExist {
		lcmmID = iavlSM.lastCommit()
	}

	iavlSApp := &IavlStoreApp{iavlSM: iavlSM, lastCommitID: lcmmID, storeLog: storeLog}

	iavlSM.storeMap[IavlStoreAccountKey].Set([]byte(AccountKey), []byte(""))
	iavlSM.storeMap[IAvlStoreMainKey].Set([]byte(CertKey), []byte(""))

	if isKVPathExist {
		iavlSApp.Prefixed(kvDB, kvPath)
	}

	return iavlSApp
}

func NewMockIavlStoreApp() *IavlStoreApp {
	db := dbm.NewMemDB()
	storeLog := log.NewNopLogger()

	iavlSM := NewIavlStoreMulti(db, storeLog)
	lcmmID := iavlSM.lastCommit()

	return  &IavlStoreApp{iavlSM: iavlSM, lastCommitID: lcmmID, storeLog: storeLog, cdc: amino.NewCodec()}
}

func (sp* IavlStoreApp) Prefixed(kvDB dbm.DB, kvPath string) error {
	var iavlStore *IavlStore
	it := kvDB.Iterator(nil, nil)

	if it != nil {
		for it.Valid() {
			if len(it.Key()) >= len(ankrtypes.AccountBlancePrefix) && string(it.Key()[0:len(ankrtypes.AccountBlancePrefix)]) == ankrtypes.AccountBlancePrefix {
				iavlStore = sp.iavlSM.IavlStore(IavlStoreAccountKey)
				keyStrList := strings.Split(string(it.Key()), ":")
				valStrList := strings.Split(string(it.Value()), ":")
				if len(keyStrList) != 2 || len(valStrList) != 2 {
					sp.storeLog.Error("invalid old account store will be ignored", "keyStrList's len", len(keyStrList), "valStrList's len", len(valStrList))
				}

				_, err := strconv.ParseInt(valStrList[1], 10, 64)
				if err == nil {
					valI, _ := new(big.Int).SetString(valStrList[0], 10)
					sp.SetBalance(keyStrList[1], account.Amount{account.Currency{"ANKR", 18}, valI.Bytes()})
				}else {
					if err != nil {
						sp.storeLog.Error("invalid old account store will be ignored: parse bal fails", "err", err)
					}
				}
			}else {
				iavlStore = sp.iavlSM.IavlStore(IAvlStoreMainKey)
				if len(it.Key()) >= len(ankrtypes.CertPrefix) && string(it.Key()[0:len(ankrtypes.CertPrefix)]) == ankrtypes.CertPrefix {
					sp.SetCertKey(it.Key(), it.Value())
				} else {
					iavlStore.Set(it.Key(), it.Value())
				}
			}
			it.Next()
		}
	}

	it.Close()
	kvDB.Close()

	err := os.RemoveAll(kvPath)

	return err
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
	if reqQuery.Path == "" {
		storeName = IAvlStoreMainKey
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
		value = []byte{} //sp.Balance(reqQuery.Data, "ANKR")
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
		value, _ = sp.AccountList()
	} else if len(reqQuery.Data) >= len(ankrtypes.AllCrtsPrefix) && string(reqQuery.Data[:len(ankrtypes.AllCrtsPrefix)]) == ankrtypes.AllCrtsPrefix {
		value = sp.CertKeyList()
	} else {
		value, _ = sp.iavlSM.IavlStore(storeName).Get(reqQuery.Data)
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

	certs, err := sp.iavlSM.IavlStore(IAvlStoreMainKey).Get([]byte(CertKey))
	if err == nil {
		certs = append(certs, []byte(";" + dcS)...)
		sp.iavlSM.IavlStore(IAvlStoreMainKey).Set([]byte(CertKey), certs)
	}else {
		sp.storeLog.Error("can't get the CertKey value", "err",  err)
	}
}

func (sp *IavlStoreApp) SetCertKey(key []byte, val []byte) {
	if !sp.iavlSM.IavlStore(IAvlStoreMainKey).Has(key) {
		dcS := strings.Split(string(key), ":")[1]
		sp.updateCertKey(dcS)
	}

	sp.iavlSM.IavlStore(IAvlStoreMainKey).Set(key, val)
}

func (sp *IavlStoreApp) CertKey(key []byte) []byte {
	valBytes, err :=  sp.iavlSM.IavlStore(IAvlStoreMainKey).Get(key)
	if err != nil {
		sp.storeLog.Error("can't get the key's value", "key", key)
		valBytes = nil
	}

	return valBytes
}

func (sp *IavlStoreApp) DeleteCertKey(key []byte) {
	sp.iavlSM.IavlStore(IAvlStoreMainKey).Remove(key)
}

func (sp *IavlStoreApp) CertKeyList() []byte {
	certs, err := sp.iavlSM.IavlStore(IAvlStoreMainKey).Get([]byte(CertKey))
	if err == nil {
		return certs
	}

	return nil
}

func (sp *IavlStoreApp) SetValidator(valInfo *ankrtypes.ValidatorInfo) {
	valBytes := ankrtypes.EncodeValidatorInfo(sp.cdc, valInfo)

	sp.iavlSM.IavlStore(IAvlStoreMainKey).Set([]byte(containValidatorPrefix(valInfo.ValAddress)), valBytes)
}

func (sp *IavlStoreApp) Validator(valAddr string) (*ankrtypes.ValidatorInfo, error) {
	valBytes, err := sp.iavlSM.IavlStore(IAvlStoreMainKey).Get([]byte(containValidatorPrefix(valAddr)))
	if err != nil {
		return nil, fmt.Errorf("can't get the responding validator info: valAddr=%s", valAddr)
	}

	valInfo := ankrtypes.DecodeValidatorInfo(sp.cdc, valBytes)

	return  &valInfo, nil
}

func (sp *IavlStoreApp) RemoveValidator(valAddr string) {
	sp.iavlSM.IavlStore(IAvlStoreMainKey).Remove([]byte(containValidatorPrefix(valAddr)))
}

func (sp *IavlStoreApp) TotalValidatorPowers() int64 {
	valPower := int64(0)

	endBytes := prefixEndBytes([]byte(StoreValidatorPrefix))

	sp.iavlSM.storeMap[IAvlStoreMainKey].tree.IterateRange([]byte(StoreValidatorPrefix), endBytes, true, func(key []byte, value []byte) bool{
		if len(key) >= len(StoreValidatorPrefix) && string(key[0:len(StoreValidatorPrefix)]) == StoreValidatorPrefix {
			valInfo := ankrtypes.DecodeValidatorInfo(sp.cdc, value)

			valPower += valInfo.Power
		}

		return false
	})

	return valPower
}

func (sp *IavlStoreApp) Get(key []byte) []byte {
	val, err := sp.iavlSM.IavlStore(IAvlStoreMainKey).Get(key)
	if err != nil {
		sp.storeLog.Error("can't get the key value", "key", string(key))
		val = nil
	}

	return val
}

func (sp *IavlStoreApp) Set(key []byte, val []byte) {
	sp.iavlSM.IavlStore(IAvlStoreMainKey).Set(key, val)
}

func (sp *IavlStoreApp) Delete(key []byte) {
	sp.iavlSM.IavlStore(IAvlStoreMainKey).Remove(key)
}

func (sp *IavlStoreApp) Has(key []byte) bool {
	return sp.iavlSM.IavlStore(IAvlStoreMainKey).Has(key)
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

func (sp *IavlStoreApp) SaveContract(key []byte, val []byte) error {
	if sp.iavlSM.IavlStore(IAvlStoreContractKey).Has(key) {
		return errors.New("the contract name has existed")
	}

	sp.iavlSM.IavlStore(IAvlStoreContractKey).Set(key, val)

	return nil
}

func (sp *IavlStoreApp) LoadContract(key []byte) ([]byte, error) {
	val, err := sp.iavlSM.IavlStore(IAvlStoreContractKey).Get(key)
	if err != nil {
		sp.storeLog.Error("can't get the contract", "key", string(key))
		val = nil
	}

	return val, err
}


