package kv

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"math/big"

	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	"github.com/tendermint/tendermint/abci/example/code"
	"github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/version"
)

var (
	stateKey        = []byte("stateKey")
	kvPairPrefixKey = []byte("kvPairKey:")

	ProtocolVersion version.Protocol = 0x1
)

type State struct {
	db      dbm.DB
	Size    int64  `json:"size"`
	Height  int64  `json:"height"`
	AppHash []byte `json:"app_hash"`
}

type Seeds struct {
    Seeds []string `json:""seeds`
}

func loadState(db dbm.DB) State {
	stateBytes := db.Get(stateKey)
	var state State
	if len(stateBytes) != 0 {
		err := json.Unmarshal(stateBytes, &state)
		if err != nil {
			panic(err)
		}
	}
	state.db = db
	return state
}

func saveState(state State) {
	stateBytes, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	state.db.Set(stateKey, stateBytes)
}

func prefixKey(key []byte) []byte {
	return append(kvPairPrefixKey, key...)
}

//---------------------------------------------------

var _ types.Application = (*KVStoreApplication)(nil)

type KVStoreApplication struct {
	types.BaseApplication

	state State
}

func NewKVStoreApplication(dbDir string) *KVStoreApplication {
	kvStore := new(KVStoreApplication)
	kvStore.init(dbDir)

	return kvStore
}

func (app *KVStoreApplication) init(dbDir string) {
	name := "kvstore"
	db, err := dbm.NewGoLevelDB(name, dbDir)
	if err != nil {
		panic(err)
	}

	app.state= loadState(db)
}

func (app *KVStoreApplication) Info(req types.RequestInfo) (resInfo types.ResponseInfo) {
	return types.ResponseInfo{
		Data:       fmt.Sprintf("{\"size\":%v}", app.state.Size),
		Version:    version.ABCIVersion,
		AppVersion: ProtocolVersion.Uint64(),
	}
}

// tx is either "key=value" or just arbitrary bytes
func (app *KVStoreApplication) DeliverTx(tx []byte) types.ResponseDeliverTx {
	var key, value []byte
	parts := bytes.Split(tx, []byte("="))
	if len(parts) == 2 {
		key, value = parts[0], parts[1]
	} else {
		key, value = tx, tx
	}
	app.state.db.Set(prefixKey(key), value)
	app.state.Size += 1

	tags := []cmn.KVPair{
		{Key: []byte("app.creator"), Value: []byte("Cosmoshi Netowoko")},
		{Key: []byte("app.key"), Value: key},
	}
	return types.ResponseDeliverTx{Code: code.CodeTypeOK, Tags: tags}
}

func (app *KVStoreApplication) CheckTx(tx []byte) types.ResponseCheckTx {
	return types.ResponseCheckTx{Code: code.CodeTypeOK, GasWanted: 1}
}

func (app *KVStoreApplication) Commit() types.ResponseCommit {
	// Using a memdb - just return the big endian size of the db
	appHash := make([]byte, 8)
	binary.PutVarint(appHash, app.state.Size)
	app.state.AppHash = appHash
	app.state.Height += 1
	saveState(app.state)
	return types.ResponseCommit{Data: appHash}
}

func (app *KVStoreApplication) Query(reqQuery types.RequestQuery) (resQuery types.ResponseQuery) {
	isBalance := false

	if reqQuery.Prove {
		value := app.state.db.Get(prefixKey(reqQuery.Data))
		resQuery.Index = -1 // TODO make Proof return index
		resQuery.Key = reqQuery.Data
		resQuery.Value = value
		if value != nil {
			resQuery.Log = "exists"
		} else {
			resQuery.Log = "does not exist"
		}
		return
	} else {
		resQuery.Key = reqQuery.Data
		value := []byte("")
		if string(reqQuery.Data[:3]) == ankrtypes.AccountBlancePrefix[:3] {
		    isBalance = true
		    value = app.state.db.Get(reqQuery.Data)
		} else if len(reqQuery.Data) >= 3 && string(reqQuery.Data[:3]) == ankrtypes.AccountStakePrefix[:3]{
		    value = app.state.db.Get(reqQuery.Data)
		} else if len(reqQuery.Data) >= 3 && string(reqQuery.Data[:3]) == ankrtypes.MeteringPrefix[:3]{
		    value = app.state.db.Get(reqQuery.Data)
		} else if len(reqQuery.Data) >= 3 && string(reqQuery.Data[:3]) == ankrtypes.CertPrefix[:3]{
		    value = app.state.db.Get(reqQuery.Data)
		} else if len(reqQuery.Data) >= 13 && string(reqQuery.Data[:13]) == ankrtypes.SET_CRT_NONCE{
		    value = app.state.db.Get(reqQuery.Data)
		} else if len(reqQuery.Data) >= 11 && string(reqQuery.Data[:11]) == ankrtypes.SET_OP_NONCE{
		    value = app.state.db.Get(reqQuery.Data)
		} else if len(reqQuery.Data) >= 9 && string(reqQuery.Data[:9]) == ankrtypes.SET_VAL_NONCE{
		    value = app.state.db.Get(reqQuery.Data)
		} else if len(reqQuery.Data) >= 13 && string(reqQuery.Data[:13]) == ankrtypes.RMV_CRT_NONCE{
		    value = app.state.db.Get(reqQuery.Data)
		} else if len(reqQuery.Data) >= 19 && string(reqQuery.Data[:19]) == ankrtypes.ADMIN_OP_VAL_PUBKEY_NAME{
		    value = app.state.db.Get(reqQuery.Data)
		} else if len(reqQuery.Data) >= 20 && string(reqQuery.Data[:20]) == ankrtypes.ADMIN_OP_FUND_PUBKEY_NAME{
		    value = app.state.db.Get(reqQuery.Data)
		} else if len(reqQuery.Data) >= 24 && string(reqQuery.Data[:24]) == ankrtypes.ADMIN_OP_METERING_PUBKEY_NAME{
		    value = app.state.db.Get(reqQuery.Data)
		} else if len(reqQuery.Data) >= len(ankrtypes.AllAccountsPrefix) &&
					string(reqQuery.Data[:len(ankrtypes.AllAccountsPrefix)]) == ankrtypes.AllAccountsPrefix {
                    itr := app.state.db.Iterator(nil, nil)
                    for ; itr.Valid(); itr.Next() {
			if len(itr.Key()) >= len(ankrtypes.AccountBlancePrefix) &&
					string(itr.Key()[0:len(ankrtypes.AccountBlancePrefix)]) == ankrtypes.AccountBlancePrefix {
			    valueItem := []byte("")
			    valueItem = app.state.db.Get(itr.Key())
			    if len(valueItem) != 0 {
				    value = []byte(string(value) + string(itr.Key()[len(ankrtypes.AccountBlancePrefix):]) + ":" + string(valueItem) + ";")
		            }
                        }
                    }
		} else if len(reqQuery.Data) >= len(ankrtypes.AllCrtsPrefix) && string(reqQuery.Data[:len(ankrtypes.AllCrtsPrefix)]) == ankrtypes.AllCrtsPrefix {
                    itr := app.state.db.Iterator(nil, nil)
                    for ; itr.Valid(); itr.Next() {
			if len(itr.Key()) >= len(ankrtypes.CertPrefix) && string(itr.Key()[0:len(ankrtypes.CertPrefix)]) == ankrtypes.CertPrefix {
			    valueItem := []byte("")
			    valueItem = app.state.db.Get(itr.Key())
			    if len(valueItem) != 0 {
				    value = []byte(string(value) + string(itr.Key()[len(ankrtypes.CertPrefix):]) + ";")
		            }
                        }
                    }
		} else if string(reqQuery.Data) == "seeds" {
            var se = []string{"127.0.0.1:26657", "127.0.0.1:26658"}
            var see = Seeds{Seeds: se}
            value, _ = json.Marshal(see)
        } else {
		    value = app.state.db.Get(prefixKey(reqQuery.Data))
		}

		//fmt.Println("queried value:", value)
		resQuery.Value = value

		if value != nil {
			if isBalance {
			    trxGetBalanceSlices := strings.Split(string(value), ":")
			    if len(trxGetBalanceSlices) == 1 {
				    _, err := new(big.Int).SetString(string(value), 10)
				    if !err {
					    resQuery.Log = "internal error, value format incorrect, single value"
					    return
				    }
			    } else if len(trxGetBalanceSlices) == 2 {
				    _, berr := new(big.Int).SetString(trxGetBalanceSlices[0], 10)
				    if !berr {
					    resQuery.Log = "internal error, value format incorrect, first value"
					    return
				    }

				    _, err := strconv.ParseInt(string(trxGetBalanceSlices[1]), 10, 64)
				    if err != nil {
					    resQuery.Log = "internal error, value format incorrect, second value"
					    return
				    }

			    } else {
				    resQuery.Log = "internal error, value format incorrect, extra value"
				    return
			    }
		        }

		        resQuery.Log = "exists"
		} else {
			resQuery.Log = "does not exist"
		}
		return
	}
}

func (app *KVStoreApplication) Get(key []byte) []byte {
	return app.state.db.Get(key)
}

func (app *KVStoreApplication) Set(key []byte, val []byte) {
	app.state.db.Set(key, val)
}

func (app *KVStoreApplication) Delete(key []byte) {
	app.state.db.Delete(key)
}

func (app *KVStoreApplication) Has(key []byte) bool {
	return app.state.db.Has(key)
}

func (app *KVStoreApplication) Validators(judgeValidatorTx ankrtypes.JudgeValidatorTx) (validators []types.ValidatorUpdate) {
	itr := app.state.db.Iterator(nil, nil)
	for ; itr.Valid(); itr.Next() {
		if judgeValidatorTx(itr.Key()) {
			validator := new(types.ValidatorUpdate)
			err := types.ReadMessage(bytes.NewBuffer(itr.Value()), validator)
			if err != nil {
				panic(err)
			}
			validators = append(validators, *validator)
		}
	}
	return
}

func (app *KVStoreApplication) TotalValidatorPowers(judgeValidatorTx ankrtypes.JudgeValidatorTx) int64 {
	var totalValPowers int64 = 0
	it := app.state.db.Iterator(nil, nil)
	if it != nil && it.Valid(){
		it.Next()
		for it.Valid() {
			if judgeValidatorTx(it.Key()) {
				validator := new(types.ValidatorUpdate)
				err := types.ReadMessage(bytes.NewBuffer(it.Value()), validator)
				if err != nil {
					panic(err)
				}

				totalValPowers += validator.Power
				fmt.Printf("validator = %v\n", validator)
			}
			it.Next()
		}
	}
	it.Close()

	return  totalValPowers
}

func (app *KVStoreApplication) Size() int64 {
	return app.state.Size
}

func (app *KVStoreApplication) IncSize() {
	app.state.Size += 1
}

func (app *KVStoreApplication) Height() int64 {
	return app.state.Height
}

func (app *KVStoreApplication) APPHash() []byte {
	return app.state.AppHash
}
