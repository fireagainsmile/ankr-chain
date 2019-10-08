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
	SetCertKey(dcName string, nsName string, pemBase64 string)
	CertKey(dcName string, nsName string) string
	DeleteCertKey(dcName string, nsName string)
	SetMetering(dcName string, nsName string, value string)
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
	IsExist(cAddr string) bool
	BuildCurrencyCAddrMap(symbol string, cAddr string)
	ContractAddrBySymbol(symbol string) (string, error)
	SaveContract(cAddr string, cInfo *ankrtypes.ContractInfo) error
	LoadContract(cAddr string) (*ankrtypes.ContractInfo, error)
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