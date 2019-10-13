package appstore

import (
	"math/big"

	ankrcmm "github.com/Ankr-network/ankr-chain/common"
	"github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
)

type AccountStore interface {
	Nonce(address string) (uint64, error)
	IncNonce(address string) (uint64, error)
	AddAccount(address string, accType ankrcmm.AccountType)
	SetBalance(address string, amount ankrcmm.Amount)
	Balance(address string, symbol string) (*big.Int, error)
	SetAllowance(addrSender string, addrSpender string, amount ankrcmm.Amount)
	Allowance(addrSender string, addrSpender string, symbol string) (*big.Int, error)
}

type BCStore interface {
	Height() int64
}

type TxStore interface {
	Commit() types.ResponseCommit
	SetCertKey(dcName string, nsName string, pemBase64 string)
	CertKey(dcName string, nsName string) string
	DeleteCertKey(dcName string, nsName string)
	SetMetering(dcName string, nsName string, value string)
	Metering(dcName string, nsName string) string
	SetValidator(valInfo *ankrcmm.ValidatorInfo)
	Validator(valAddr string) (*ankrcmm.ValidatorInfo, error)
	RemoveValidator(valAddr string)
	TotalValidatorPowers() int64
	Get(key []byte) []byte
	Set(key []byte, val []byte)
	Delete(key []byte)
	Has(key []byte) bool
}

type ContractStore interface {
	IsExist(cAddr string) bool
	BuildCurrencyCAddrMap(symbol string, cAddr string) error
	ContractAddrBySymbol(symbol string) (string, error)
	SaveContract(cAddr string, cInfo *ankrcmm.ContractInfo) error
	LoadContract(cAddr string) (*ankrcmm.ContractInfo, error)
}

type QueryHandler interface {
	Query(reqQuery types.RequestQuery) (resQuery types.ResponseQuery)
}

type AppStore interface {
	AccountStore
	TxStore
	QueryHandler
	ContractStore
	BCStore
	APPHash() []byte
    DB() dbm.DB
}