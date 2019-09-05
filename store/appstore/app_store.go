package appstore

import (
	"github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
)

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

type ContractStore interface {
	SaveContract(key []byte, val []byte) error
	LoadContract(key []byte) ([]byte, error)
}

type QueryHandler interface {
	Query(reqQuery types.RequestQuery) (resQuery types.ResponseQuery)
}

type AppStore interface {
	AccountStore
	TxStore
	QueryHandler
	ContractStore
	Height() int64
	APPHash() []byte
    DB() dbm.DB
}