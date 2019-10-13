package iavl

import (
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"

	ankrcmm "github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/common/code"
	apscomm "github.com/Ankr-network/ankr-chain/store/appstore/common"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/merkle"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

const (
	StoreCertKeyPrefix = "certkey:"
	StoreMeteringPrefix = "merting:"
	StoreContractInfoPrefix = "continfo:"
	StoreContractCurrencyPrefix = "contcur:"
)

type storeQueryHandler func(store *IavlStoreApp, reqData []byte) (resQuery types.ResponseQuery)

type IavlStoreApp struct {
	iavlSM          *IavlStoreMulti
	lastCommitID    ankrcmm.CommitID
	storeLog        log.Logger
	cdc             *amino.Codec
	queryHandleMap map[string] storeQueryHandler
	accStoreLocker  sync.RWMutex
	certStoreLocker sync.RWMutex
}

func containCertKeyPrefix(dcnsName string) string {
	return containPrefix(dcnsName, StoreCertKeyPrefix)
}

func containMeteringPrefix(dcnsName string) string {
	return containPrefix(dcnsName, StoreMeteringPrefix)
}

func containContractInfoPrefix(cAddr string) string {
	return containPrefix(cAddr, StoreContractInfoPrefix)
}

func containContractCurrencyPrefix(symbol string) string {
	return containPrefix(symbol, StoreContractCurrencyPrefix)
}

func stripCertKeyPrefix(key string) (string, error) {
	return stripKeyPrefix(key, StoreCertKeyPrefix)
}

func NewIavlStoreApp(dbDir string, storeLog log.Logger) *IavlStoreApp {
	kvPath := filepath.Join(dbDir, "kvstore.db")
	isKVPathExist, err := ankrcmm.PathExists(kvPath)
	if err != nil {
		panic(err)
	}

	var kvDB dbm.DB
	var lcmmID ankrcmm.CommitID
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

	if isKVPathExist {
		iavlSApp.Prefixed(kvDB, kvPath)
	}

	iavlSApp.queryHandleMap = make(map[string]storeQueryHandler)

	iavlSApp.registerQueryHandlerWapper("balance",   ankrcmm.BalanceQueryReq{},  iavlSApp.Balance)
	iavlSApp.registerQueryHandlerWapper("certkey",   ankrcmm.CertKeyQueryReq{},  iavlSApp.CertKey)
	iavlSApp.registerQueryHandlerWapper("metering",  ankrcmm.MeteringQueryReq{}, iavlSApp.Metering)
	iavlSApp.registerQueryHandlerWapper("validator", ankrcmm.ValidatorQueryReq{},iavlSApp.Validator)
	iavlSApp.registerQueryHandlerWapper("contract",  ankrcmm.ContractQueryReq{}, iavlSApp.LoadContract)

	return iavlSApp
}

func NewMockIavlStoreApp() *IavlStoreApp {
	db := dbm.NewMemDB()
	storeLog := log.NewNopLogger()

	iavlSM := NewIavlStoreMulti(db, storeLog)
	lcmmID := iavlSM.lastCommit()

	iavlSApp := &IavlStoreApp{iavlSM: iavlSM, lastCommitID: lcmmID, storeLog: storeLog, cdc: amino.NewCodec()}

	iavlSApp.queryHandleMap = make(map[string]storeQueryHandler)

	iavlSApp.registerQueryHandlerWapper("balance",   ankrcmm.BalanceQueryReq{},  iavlSApp.Balance)
	iavlSApp.registerQueryHandlerWapper("certkey",   ankrcmm.CertKeyQueryReq{},  iavlSApp.CertKey)
	iavlSApp.registerQueryHandlerWapper("metering",  ankrcmm.MeteringQueryReq{}, iavlSApp.Metering)
	iavlSApp.registerQueryHandlerWapper("validator", ankrcmm.ValidatorQueryReq{},iavlSApp.Validator)
	iavlSApp.registerQueryHandlerWapper("contract",  ankrcmm.ContractQueryReq{}, iavlSApp.LoadContract)

	return  &IavlStoreApp{iavlSM: iavlSM, lastCommitID: lcmmID, storeLog: storeLog, cdc: amino.NewCodec()}
}


func (sp* IavlStoreApp) registerQueryHandlerWapper(queryKey string, req interface{}, callFunc interface{}) {
	sp.queryHandleMap[queryKey] = func(store *IavlStoreApp, reqData []byte) (resQuery types.ResponseQuery) {
		err := store.cdc.UnmarshalJSON(reqData, req)
		if err != nil {
			resQuery.Code = code.CodeTypeQueryInvalidQueryReqData
			resQuery.Log  = fmt.Sprintf("invalid %s query req data, err=%s", queryKey, err.Error())
			return
		}

		v := reflect.ValueOf(callFunc)
		var paramVals [] reflect.Value
		reqVals := reflect.ValueOf(req)
		for i := 0; i < reqVals.NumField(); i++ {
			paramVals = append(paramVals, reqVals.Field(i))
		}
		respVals := v.Call(paramVals)

		for j := 0; j < len(respVals); j++ {
			if respVals[j].Interface() != nil && respVals[j].Type().Name() == reflect.TypeOf(errors.New("")).Name() {
				err :=  respVals[j].Interface().(error)
				resQuery.Code = code.CodeTypeLoadBalError
				resQuery.Log  = fmt.Sprintf("load %s query err, err=%s", queryKey, err.Error())
				return
			}

			if respVals[j].Interface() != nil {
				resQDataBytes, _ := store.cdc.MarshalJSON(respVals[0].Interface())
				resQuery.Code  = code.CodeTypeOK
				resQuery.Value = resQDataBytes
			}

		}

		resQuery.Code  = code.CodeTypeUnknownError

		return
	}
}

func (sp* IavlStoreApp) Prefixed(kvDB dbm.DB, kvPath string) error {
	var iavlStore *IavlStore
	it := kvDB.Iterator(nil, nil)

	if it != nil {
		for it.Valid() {
			if len(it.Key()) >= len(ankrcmm.AccountBlancePrefix) && string(it.Key()[0:len(ankrcmm.AccountBlancePrefix)]) == ankrcmm.AccountBlancePrefix {
				iavlStore = sp.iavlSM.IavlStore(IavlStoreAccountKey)
				keyStrList := strings.Split(string(it.Key()), ":")
				valStrList := strings.Split(string(it.Value()), ":")
				if len(keyStrList) != 2 || len(valStrList) != 2 {
					sp.storeLog.Error("invalid old account store will be ignored", "keyStrList's len", len(keyStrList), "valStrList's len", len(valStrList))
				}

				_, err := strconv.ParseInt(valStrList[1], 10, 64)
				if err == nil {
					valI, _ := new(big.Int).SetString(valStrList[0], 10)
					sp.SetBalance(keyStrList[1], ankrcmm.Amount{ankrcmm.Currency{"ANKR", 18}, valI.Bytes()})
				}else {
					if err != nil {
						sp.storeLog.Error("invalid old account store will be ignored: parse bal fails", "err", err)
					}
				}
			}else {
				iavlStore = sp.iavlSM.IavlStore(IAvlStoreMainKey)
				if len(it.Key()) >= len(ankrcmm.CertPrefix) && string(it.Key()[0:len(ankrcmm.CertPrefix)]) == ankrcmm.CertPrefix {
					//sp.SetCertKey(it.Key(), it.Value())
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

func (sp *IavlStoreApp) LastCommit() *ankrcmm.CommitID{
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
	resQuery.Log = "exists"
	if reqQuery.Path == "" {
		resQuery.Log = "blank store name"
		resQuery.Code = code.CodeTypeQueryInvalidStoreName
		return
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
				Key:  []byte(reqQuery.Path),
				Data: infoBytes,
			}
			resQuery.Proof.Ops = append(resQuery.Proof.Ops, pOP)
		}
	}

    return sp.queryHandleMap[reqQuery.Path](sp, reqQuery.Data)
}

func (sp *IavlStoreApp) SetCertKey(dcName string, nsName string, pemBase64 string)  {
	key := []byte(containCertKeyPrefix(dcName+"_"+nsName))
	sp.iavlSM.IavlStore(IAvlStoreMainKey).Set(key, []byte(pemBase64))
}

func (sp *IavlStoreApp) CertKey(dcName string, nsName string) string {
	key := []byte(containCertKeyPrefix(dcName+"_"+nsName))
	valBytes, err :=  sp.iavlSM.IavlStore(IAvlStoreMainKey).Get(key)
	if err != nil {
		sp.storeLog.Error("can't get the key's value", "dcName", dcName)
		return ""
	}

	return string(valBytes)
}

func (sp *IavlStoreApp) DeleteCertKey(dcName string, nsName string) {
	key := []byte(containCertKeyPrefix(dcName+"_"+nsName))
	sp.iavlSM.IavlStore(IAvlStoreMainKey).Remove(key)
}

func (sp *IavlStoreApp) CertKeyList() ([]byte, uint64) {
	certKeyCount := uint64(0)
	certKeyList  := ""

	endBytes := prefixEndBytes([]byte(StoreCertKeyPrefix))

	sp.iavlSM.storeMap[IAvlStoreMainKey].tree.IterateRange([]byte(StoreCertKeyPrefix), endBytes, true, func(key []byte, value []byte) bool{
		if len(key) >= len(StoreAccountPrefix) && string(key[0:len(StoreCertKeyPrefix)]) == StoreCertKeyPrefix {
			dcnsName, err := stripCertKeyPrefix(string(key))
			if err != nil {
				sp.storeLog.Error("stripCertKeyPrefix error", "err", err)
			}else {
				certKeyCount++
				certKeyList += certKeyList + ";" + dcnsName + ":" + string(value)
			}
		}

		return false
	})

	if certKeyCount > 0 {
		certKeyList = certKeyList[1:]
		return []byte(certKeyList), certKeyCount
	}else {
		return nil, certKeyCount
	}
}

func (sp *IavlStoreApp) SetMetering(dcName string, nsName string, value string) {
	key := []byte(containCertKeyPrefix(dcName+"_"+nsName))
	sp.iavlSM.IavlStore(IAvlStoreMainKey).Set(key, []byte(value))
}

func (sp *IavlStoreApp) Metering(dcName string, nsName string) string {
	key := []byte(containCertKeyPrefix(dcName+"_"+nsName))
	valueBytes, err := sp.iavlSM.IavlStore(IAvlStoreMainKey).Get(key)
	if err != nil {
		sp.storeLog.Error("can't get the responding metering value", "dcName", dcName, "nsName", nsName, "err", err)
		return ""
	}

	return string(valueBytes)
}

func (sp *IavlStoreApp) SetValidator(valInfo *ankrcmm.ValidatorInfo) {
	valBytes := ankrcmm.EncodeValidatorInfo(sp.cdc, valInfo)

	sp.iavlSM.IavlStore(IAvlStoreMainKey).Set([]byte(containValidatorPrefix(valInfo.ValAddress)), valBytes)
}

func (sp *IavlStoreApp) Validator(valAddr string) (*ankrcmm.ValidatorInfo, error) {
	valBytes, err := sp.iavlSM.IavlStore(IAvlStoreMainKey).Get([]byte(containValidatorPrefix(valAddr)))
	if err != nil {
		return nil, fmt.Errorf("can't get the responding validator info: valAddr=%s", valAddr)
	}

	valInfo := ankrcmm.DecodeValidatorInfo(sp.cdc, valBytes)

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
			valInfo := ankrcmm.DecodeValidatorInfo(sp.cdc, value)

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

func (sp *IavlStoreApp) IsExist(cAddr string) bool {
	return sp.iavlSM.IavlStore(IAvlStoreContractKey).Has([]byte(cAddr))
}

func (sp *IavlStoreApp) BuildCurrencyCAddrMap(symbol string, cAddr string) error {
	if sp.iavlSM.IavlStore(IAvlStoreContractKey).Has([]byte(containContractCurrencyPrefix(symbol))) {
		return errors.New("the contract name has existed")
	}

	sp.iavlSM.IavlStore(IAvlStoreContractKey).Set([]byte(containContractCurrencyPrefix(symbol)), []byte(cAddr))

	return nil
}

func (sp *IavlStoreApp) ContractAddrBySymbol(symbol string) (string, error) {
	cAddrBytes, err := sp.iavlSM.IavlStore(IAvlStoreContractKey).Get([]byte(containContractCurrencyPrefix(symbol)))
	if err != nil {
		sp.storeLog.Error("can't get the contract addr", "symbol", symbol)
		return "", err
	}

	if cAddrBytes != nil {
		return string(cAddrBytes), nil
	}

	return "", nil

}

func (sp *IavlStoreApp) SaveContract(cAddr string, cInfo *ankrcmm.ContractInfo) error{
	if sp.iavlSM.IavlStore(IAvlStoreContractKey).Has([]byte(containContractInfoPrefix(cAddr))) {
		return errors.New("the contract name has existed")
	}

	cInfoBytes := ankrcmm.EncodeContractInfo(sp.cdc, cInfo)

	sp.iavlSM.IavlStore(IAvlStoreContractKey).Set([]byte(containContractInfoPrefix(cAddr)), cInfoBytes)

	return nil
}

func (sp *IavlStoreApp) LoadContract(cAddr string) (*ankrcmm.ContractInfo, error) {
	cInfoBytes, err := sp.iavlSM.IavlStore(IAvlStoreContractKey).Get([]byte(containContractInfoPrefix(cAddr)))
	if err != nil {
		sp.storeLog.Error("can't get the contract", "addr", cAddr)
		return nil, err
	}

	cInfo := ankrcmm.DecodeContractInfo(sp.cdc, cInfoBytes)

	return &cInfo, nil
}


