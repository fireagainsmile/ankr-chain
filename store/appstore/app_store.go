package appstore

import (
	"math/big"

	"github.com/Ankr-network/ankr-chain/account"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	"github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"

)

type AccountStore interface {
	InitGenesisAccount()
	InitFoundAccount()
	Nonce(address string) (uint64, error)
	IncNonce(address string) (uint64, error)
	SetBalance(address string, amount account.Amount)
	Balance(address string, symbol string) (*big.Int, error)
	SetAllowance(addrSender string, addrSpender string, amount account.Amount)
	Allowance(addrSender string, addrSpender string, symbol string) (*big.Int, error)
}

type TxStore interface {
	Commit() types.ResponseCommit
	SetCertKey(key []byte, val []byte)
	CertKey(key []byte) []byte
	DeleteCertKey(key []byte)
	SetValidator(valInfo *ankrtypes.ValidatorInfo)
	Validator(valAddr string) (*ankrtypes.ValidatorInfo, error)
	RemoveValidator(valAddr string)
	TotalValidatorPowers() int64
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