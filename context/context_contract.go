package context

import (
	"math/big"
	"context"

	"github.com/Ankr-network/ankr-chain/account"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	vmevent "github.com/go-interpreter/wagon/exec/event"
	"github.com/go-interpreter/wagon/exec/gas"
)

type TxMsgCallBack interface {
	SenderAddr() string
}

type ContextContract interface {
	SpendGas(gas *big.Int) bool
	SenderAddr() string
	OwnerAddr() string
	ContractAddr() string
	SetBalance(address string, amount account.Amount)
	Balance(address string, symbol string) (*big.Int, error)
	SetAllowance(addrSender string, addrSpender string, amount account.Amount)
	Allowance(addrSender string, addrSpender string, symbol string) (*big.Int, error)
	Publish(ctx context.Context, msg interface{}) error
	PublishWithTags(ctx context.Context, msg interface{}, tags map[string]string) error
}

type ContextContractImpl struct {
	gas.GasMetric
	TxMsgCallBack
	ankrtypes.ContractInterface
	appstore.AccountStore
	vmevent.Publisher
}

func NewContextContract(gasMetric gas.GasMetric, txCallBack TxMsgCallBack, contI ankrtypes.ContractInterface, accStore appstore.AccountStore, publisher vmevent.Publisher) ContextContract {
	return &ContextContractImpl{gasMetric, txCallBack, contI,accStore, publisher}
}


