package context

import (
	"context"
	"math/big"

	ankrcmm "github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	vmevent "github.com/go-interpreter/wagon/exec/event"
	"github.com/go-interpreter/wagon/exec/gas"
)

type TxMsgCallBack interface {
	SenderAddr() string
}

type CurrencyInterface interface {
	CreateCurrency(symbol string, currency *ankrcmm.Currency) error
	CurrencyInfo(symbol string) (*ankrcmm.Currency, error)
}

type ContextContract interface {
	CreateCurrency(symbol string, currency *ankrcmm.Currency) error
	CurrencyInfo(symbol string) (*ankrcmm.Currency, error)
	SpendGas(gas *big.Int) bool
	SenderAddr() string
	OwnerAddr() string
	ContractAddr() string
	SetBalance(address string, amount ankrcmm.Amount)
	Balance(address string, symbol string) (*big.Int, error)
	SetAllowance(addrSender string, addrSpender string, amount ankrcmm.Amount)
	Allowance(addrSender string, addrSpender string, symbol string) (*big.Int, error)
	Publish(ctx context.Context, msg interface{}) error
	PublishWithTags(ctx context.Context, msg interface{}, tags map[string]string) error
}

type ContextContractImpl struct {
	CurrencyInterface
	gas.GasMetric
	TxMsgCallBack
	ankrcmm.ContractInterface
	appstore.AccountStore
	vmevent.Publisher
}

func NewContextContract(curI CurrencyInterface, gasMetric gas.GasMetric, txCallBack TxMsgCallBack, contI ankrcmm.ContractInterface, accStore appstore.AccountStore, publisher vmevent.Publisher) ContextContract {
	return &ContextContractImpl{curI,gasMetric, txCallBack, contI,accStore, publisher}
}


