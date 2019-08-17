package kv

import (
	"encoding/binary"
	"encoding/json"
	"math/big"
	"strconv"
	"strings"

	"github.com/Ankr-network/ankr-chain/store/appstore"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	"github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/version"
)

var (
	kvPairPrefixKey = []byte("kvPairKey:")

	ProtocolVersion version.Protocol = 0x1
)

func saveState(state appstore.State) {
	stateBytes, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	state.DB.Set(appstore.StateKey, stateBytes)
}

func prefixKey(key []byte) []byte {
	return append(kvPairPrefixKey, key...)
}

var _ types.Application = (*KVStoreApplication)(nil)

type KVStoreApplication struct {
	types.BaseApplication

	state appstore.State
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

	app.state= appstore.LoadState(db)
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
	var value []byte
	resQuery.Log = "exists"

	if reqQuery.Prove {
		value = app.state.DB.Get(prefixKey(reqQuery.Data))
		resQuery.Index = -1 // TODO make Proof return index
		resQuery.Key = reqQuery.Data
		resQuery.Value = value
	} else {
		resQuery.Key = reqQuery.Data
		if string(reqQuery.Data[:3]) == ankrtypes.AccountBlancePrefix[:3] {
		    value = app.state.DB.Get(reqQuery.Data)
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
			itr := app.state.DB.Iterator(nil, nil)
			for ; itr.Valid(); itr.Next() {
				if len(itr.Key()) >= len(ankrtypes.AccountBlancePrefix) && string(itr.Key()[0:len(ankrtypes.AccountBlancePrefix)]) == ankrtypes.AccountBlancePrefix {
					valueItem := []byte("")
					valueItem = app.state.DB.Get(itr.Key())
					if len(valueItem) != 0 {
						value = []byte(string(value) + string(itr.Key()[len(ankrtypes.AccountBlancePrefix):]) + ":" + string(valueItem) + ";")
					 }
				}
			}
		} else if len(reqQuery.Data) >= len(ankrtypes.AllCrtsPrefix) && string(reqQuery.Data[:len(ankrtypes.AllCrtsPrefix)]) == ankrtypes.AllCrtsPrefix {
			itr := app.state.DB.Iterator(nil, nil)
			for ; itr.Valid(); itr.Next() {
				if len(itr.Key()) >= len(ankrtypes.CertPrefix) && string(itr.Key()[0:len(ankrtypes.CertPrefix)]) == ankrtypes.CertPrefix {
					valueItem := []byte("")
					valueItem = app.state.DB.Get(itr.Key())
					if len(valueItem) != 0 {
						value = []byte(string(value) + string(itr.Key()[len(ankrtypes.CertPrefix):]) + ";")
					 }
				}
            }
		} else {
			value = app.state.DB.Get(reqQuery.Data)
		}
	}

	resQuery.Value = value

	if value == nil {
		resQuery.Log = "does not exist"
	}

	return
}

func (app *KVStoreApplication) SetCertKey(key []byte, val []byte) {
	app.state.DB.Set(key, val)
}

func (app *KVStoreApplication) CertKey(key []byte) []byte {
	return app.state.DB.Get(key)
}

func (app *KVStoreApplication) DeleteCertKey(key []byte) {
	app.state.DB.Delete(key)
}

func (app *KVStoreApplication) Get(key []byte) []byte {
	return app.state.DB.Get(key)
}

func (app *KVStoreApplication) Set(key []byte, val []byte) {
	app.state.DB.Set(key, val)
}

func (app *KVStoreApplication) Delete(key []byte) {
	app.state.DB.Delete(key)
}

func (app *KVStoreApplication) Has(key []byte) bool {
	return app.state.DB.Has(key)
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

func (app *KVStoreApplication) DB() dbm.DB {
	return app.state.DB
}

func (app *KVStoreApplication) SetBalance(key []byte, val []byte) {
	app.state.DB.Set(key, val)
}

func (app *KVStoreApplication) Balance(key []byte) []byte {
	return app.state.DB.Get(key)
}
