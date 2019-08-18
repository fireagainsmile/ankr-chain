package appstore

import (
	"encoding/json"

	"github.com/Ankr-network/ankr-chain/router"
	"github.com/Ankr-network/ankr-chain/store/appstore/iavl"
	"github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

var (
	StateKey = []byte("stateKey")
)

type State struct {
	DB      dbm.DB
	Size    int64  `json:"size"`
	Height  int64  `json:"height"`
	AppHash []byte `json:"app_hash"`
}

type AccountStore interface {
	SetBalance(key []byte, val []byte)
	Balance(key []byte) []byte
}

type TxStore interface {
	Commit() types.ResponseCommit
	SetCertKey(key []byte, val []byte)
	CertKey(key []byte) []byte
	DeleteCertKey(key []byte)
	Get(key []byte) []byte
	Set(key []byte, val []byte)
	Delete(key []byte)
	Has(key []byte) bool
}

type AppStore interface {
	AccountStore
	TxStore
	router.QueryHandler
	Height() int64
	APPHash() []byte
    DB() dbm.DB
}

func LoadState(db dbm.DB) State {
	stateBytes := db.Get(StateKey)
	var state State
	if len(stateBytes) != 0 {
		err := json.Unmarshal(stateBytes, &state)
		if err != nil {
			panic(err)
		}
	}
	state.DB = db

	return state
}

func NewAppStore(dbDir string, l log.Logger) AppStore {
	appStore := iavl.NewIavlStoreApp(dbDir, l)

	router.QueryRouterInstance().AddQueryHandler("store", appStore)

	return  appStore
}