package context

import (
	"github.com/Ankr-network/ankr-chain/tx"
	"math/big"

	"github.com/Ankr-network/ankr-chain/account"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	"github.com/go-interpreter/wagon/gas"
)

type ContextContract interface {
	SpendGas(gas *big.Int) bool
	SenderAddr() string
	SetBalance(address string, amount account.Amount)
	Balance(address string, symbol string) (*big.Int, error)
	SetAllowance(addrSender string, addrSpender string, amount account.Amount)
	Allowance(addrSender string, addrSpender string, symbol string) (*big.Int, error)
}

type ContextContractImpl struct {
	gas.GasMetric
	tx.TxMsgCallBack
	appstore.AccountStore
}

func NewContextContract(gasMetric gas.GasMetric, txCallBack tx.TxMsgCallBack, accStore appstore.AccountStore) ContextContract {
	return &ContextContractImpl{gasMetric, txCallBack, accStore}
}


