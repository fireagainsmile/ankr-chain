package iavl

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	ankrcmm "github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/common/code"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/merkle"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

const (
	ChainIDKey = "chainidkey"
	TotalTxKey = "totaltxkey"
)

const (
	StoreCertKeyPrefix = "certkey:"
	StoreMeteringPrefix = "merting:"
	StoreContractInfoPrefix = "continfo:"
	StoreContractCurrencyPrefix = "contcur:"
)

//type storeQueryHandler func(store *IavlStoreApp, reqData []byte) (resQuery types.ResponseQuery)

type IavlStoreApp struct {
	iavlSM          *IavlStoreMulti
	lastCommitID    ankrcmm.CommitID
	totalTx         int64
	storeLog        log.Logger
	cdc             *amino.Codec
	queryHandleMap  map[string]*storeQueryHandler
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

type storeQueryHandler struct {
	req interface{}
	callFunc interface{}
}

func NewIavlStoreApp(dbDir string, storeLog log.Logger) *IavlStoreApp {
	var lcmmID ankrcmm.CommitID

	db, err := dbm.NewGoLevelDB("appstore", dbDir)
	if err != nil {
		panic(err)
	}

	iavlSM := NewIavlStoreMulti(db, storeLog)

	lcmmID = iavlSM.lastCommit()

	iavlSM.Load()

	iavlSApp := &IavlStoreApp{iavlSM: iavlSM, lastCommitID: lcmmID, storeLog: storeLog, cdc: amino.NewCodec()}

	if !iavlSM.storeMap[IAvlStoreMainKey].Has([]byte(TotalTxKey)) {
		buf := make([]byte, binary.MaxVarintLen64)
		n := binary.PutVarint(buf, int64(0))
		iavlSM.storeMap[IAvlStoreMainKey].Set([]byte(TotalTxKey), buf[:n])
	} else {
		totalTxBytes, err := iavlSM.storeMap[IAvlStoreMainKey].Get([]byte(TotalTxKey))
		if err == nil {
			iavlSApp.totalTx, _ = binary.Varint(totalTxBytes)
		}else {
			storeLog.Error("load txtal tx error", "err", err)
		}
	}

	lastHash := make([]byte, 8)
	binary.PutVarint(lastHash, iavlSApp.totalTx)

	if lcmmID.Hash != nil {
		lcmmID.Hash = lastHash
	}

	iavlSApp.lastCommitID = lcmmID

	iavlSApp.queryHandleMap = make(map[string]*storeQueryHandler)

	iavlSApp.queryHandleMap["nonce"]     = &storeQueryHandler{&ankrcmm.NonceQueryReq{},  iavlSApp.NonceQuery}
	iavlSApp.queryHandleMap["balance"]   = &storeQueryHandler{&ankrcmm.BalanceQueryReq{},  iavlSApp.BalanceQuery}
	iavlSApp.queryHandleMap["certkey"]   = &storeQueryHandler{&ankrcmm.CertKeyQueryReq{},  iavlSApp.CertKeyQuery}
	iavlSApp.queryHandleMap["metering"]  = &storeQueryHandler{&ankrcmm.MeteringQueryReq{}, iavlSApp.MeteringQuery}
	iavlSApp.queryHandleMap["validator"] = &storeQueryHandler{&ankrcmm.ValidatorQueryReq{},iavlSApp.ValidatorQuery}
	iavlSApp.queryHandleMap["contract"]  = &storeQueryHandler{&ankrcmm.ContractQueryReq{}, iavlSApp.LoadContractQuery}

	return iavlSApp
}

func NewMockIavlStoreApp() *IavlStoreApp {
	db := dbm.NewMemDB()
	storeLog := log.NewNopLogger()

	iavlSM := NewIavlStoreMulti(db, storeLog)
	lcmmID := iavlSM.lastCommit()

	iavlSApp := &IavlStoreApp{iavlSM: iavlSM, lastCommitID: lcmmID, storeLog: storeLog, cdc: amino.NewCodec()}

	iavlSApp.queryHandleMap = make(map[string]*storeQueryHandler)

	iavlSApp.queryHandleMap["nonce"]     = &storeQueryHandler{&ankrcmm.NonceQueryReq{},    iavlSApp.NonceQuery}
	iavlSApp.queryHandleMap["balance"]   = &storeQueryHandler{&ankrcmm.BalanceQueryReq{},  iavlSApp.BalanceQuery}
	iavlSApp.queryHandleMap["certkey"]   = &storeQueryHandler{&ankrcmm.CertKeyQueryReq{},  iavlSApp.CertKeyQuery}
	iavlSApp.queryHandleMap["metering"]  = &storeQueryHandler{&ankrcmm.MeteringQueryReq{}, iavlSApp.MeteringQuery}
	iavlSApp.queryHandleMap["validator"] = &storeQueryHandler{&ankrcmm.ValidatorQueryReq{},iavlSApp.ValidatorQuery}
	iavlSApp.queryHandleMap["contract"]  = &storeQueryHandler{&ankrcmm.ContractQueryReq{}, iavlSApp.LoadContractQuery}

	return  &IavlStoreApp{iavlSM: iavlSM, lastCommitID: lcmmID, storeLog: storeLog, cdc: amino.NewCodec()}
}

func (sp* IavlStoreApp) queryHandlerWapper(queryKey string, reqData []byte) (resQuery types.ResponseQuery) {
	req      := sp.queryHandleMap[queryKey].req
	callFunc := sp.queryHandleMap[queryKey].callFunc
	err := sp.cdc.UnmarshalJSON(reqData, req)
	if err != nil {
		resQuery.Code = code.CodeTypeQueryInvalidQueryReqData
		resQuery.Log  = fmt.Sprintf("invalid %s query req data, err=%s", queryKey, err.Error())
		return
	}

	v := reflect.ValueOf(callFunc)
	var paramVals [] reflect.Value
	reqVals := reflect.ValueOf(req)
	reqVals = reflect.Indirect(reqVals)
	for i := 0; i < reqVals.NumField(); i++ {
		paramVals = append(paramVals, reqVals.Field(i))
	}
	respVals := v.Call(paramVals)

	if respVals[1].Interface() != nil && respVals[1].Type().Name() == reflect.TypeOf(errors.New("")).Name() {
		err :=  respVals[1].Interface().(error)
		resQuery.Code = code.CodeTypeLoadBalError
		resQuery.Log  = fmt.Sprintf("load %s query err, err=%s", queryKey, err.Error())
		return
	}

	if respVals[1].Interface() == nil {
		resQDataBytes, _ := sp.cdc.MarshalJSON(respVals[0].Interface())
		resQuery.Code    = code.CodeTypeOK
		resQuery.Value   = resQDataBytes
		return
	}

	resQuery.Code  = code.CodeTypeUnknownError
	resQuery.Log   = fmt.Sprintf("load %s query unknown err", queryKey)

	return
}

func (sp *IavlStoreApp) SetChainID(chainID string) {
	sp.iavlSM.IavlStore(IAvlStoreMainKey).Set([]byte(ChainIDKey), []byte(chainID))
}

func (sp *IavlStoreApp) ChainID() string {
	chainIDBytes, err := sp.iavlSM.IavlStore(IAvlStoreMainKey).Get([]byte(ChainIDKey))
	if err != nil || chainIDBytes == nil{
		return ""
	}

	return string(chainIDBytes)
}

func (sp *IavlStoreApp) LastCommit() *ankrcmm.CommitID{
	return &sp.lastCommitID
}

func (sp *IavlStoreApp) Commit() types.ResponseCommit {
    commitID := sp.iavlSM.Commit(sp.lastCommitID.Version)

	appHash := make([]byte, 8)
	binary.PutVarint(appHash, sp.totalTx)

	sp.lastCommitID.Hash = sp.lastCommitID.Hash[0:0]

    sp.lastCommitID.Version = commitID.Version
	sp.lastCommitID.Hash    = append(sp.lastCommitID.Hash, appHash...)

	return types.ResponseCommit{Data: appHash}
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

    return sp.queryHandlerWapper(reqQuery.Path, reqQuery.Data)
}

func (sp *IavlStoreApp) SetCertKey(dcName string, pemBase64 string)  {
	key := []byte(containCertKeyPrefix(dcName))
	sp.iavlSM.IavlStore(IAvlStoreMainKey).Set(key, []byte(pemBase64))
}

func (sp *IavlStoreApp) CertKey(dcName string) string {
	key := []byte(containCertKeyPrefix(dcName))
	valBytes, err :=  sp.iavlSM.IavlStore(IAvlStoreMainKey).Get(key)
	if err != nil {
		sp.storeLog.Error("can't get the key's value", "dcName", dcName)
		return ""
	}

	return string(valBytes)
}

func (sp *IavlStoreApp) CertKeyQuery(dcName string) (*ankrcmm.CertKeyQueryResp, error) {
	return &ankrcmm.CertKeyQueryResp{sp.CertKey(dcName)}, nil
}

func (sp *IavlStoreApp) DeleteCertKey(dcName string) {
	key := []byte(containCertKeyPrefix(dcName))
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
	key := []byte(containMeteringPrefix(dcName+"_"+nsName))
	sp.iavlSM.IavlStore(IAvlStoreMainKey).Set(key, []byte(value))
}

func (sp *IavlStoreApp) Metering(dcName string, nsName string) string {
	key := []byte(containMeteringPrefix(dcName+"_"+nsName))
	valueBytes, err := sp.iavlSM.IavlStore(IAvlStoreMainKey).Get(key)
	if err != nil {
		sp.storeLog.Error("can't get the responding metering value", "dcName", dcName, "nsName", nsName, "err", err)
		return ""
	}

	return string(valueBytes)
}

func (sp *IavlStoreApp) MeteringQuery(dcName string, nsName string) (*ankrcmm.MeteringQueryResp, error) {
	return &ankrcmm.MeteringQueryResp{sp.Metering(dcName, nsName)}, nil
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

func (sp *IavlStoreApp) ValidatorQuery(valAddr string) (*ankrcmm.ValidatorQueryResp, error) {
	valInfo, err := sp.Validator(valAddr)
	if err != nil {
		return nil, err
	}

	valQResp := &ankrcmm.ValidatorQueryResp{
		Name:         valInfo.Name,
		ValAddress:   valInfo.ValAddress,
		PubKey:       valInfo.PubKey,
		Power:        valInfo.Power,
		StakeAddress: valInfo.StakeAddress,
		StakeAmount:  valInfo.StakeAmount,
		ValidHeight:  valInfo.ValidHeight,
	}

	return valQResp, nil
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

func (sp *IavlStoreApp) TotalTx() int64 {
	return sp.totalTx
}

func (sp *IavlStoreApp) IncTotalTx() int64 {
	sp.totalTx++

	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(buf, sp.totalTx)
	sp.iavlSM.storeMap[IAvlStoreMainKey].Set([]byte(TotalTxKey), buf[:n])

	return sp.totalTx
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
	if err != nil || len(cInfoBytes) == 0{
		sp.storeLog.Error("can't get the contract", "addr", cAddr)
		return nil, err
	}

	cInfo := ankrcmm.DecodeContractInfo(sp.cdc, cInfoBytes)

	return &cInfo, nil
}

func (sp *IavlStoreApp) LoadContractQuery(cAddr string) (*ankrcmm.ContractQueryResp, error) {
	cInfo, err := sp.LoadContract(cAddr)
	if err != nil {
		return nil, err
	}

	cInfoQResp := &ankrcmm.ContractQueryResp{}
	cInfoQResp.Addr  = cInfo.Addr
	cInfoQResp.Name  = cInfo.Name
	cInfoQResp.Owner = cInfo.Owner

	cInfoQResp.Codes = make([]byte, len(cInfo.Codes))
	copy(cInfoQResp.Codes[:], cInfo.Codes)

	cInfoQResp.CodesDesc = cInfo.CodesDesc

	return cInfoQResp, nil
}


