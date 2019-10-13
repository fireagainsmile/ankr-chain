package context

import (
	"context"
	"math/big"

	ankrcmm "github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/store/appstore"
)

var bcContext ContextAKVM

type ContextAKVM interface {
	SpendGas(gas *big.Int) bool
	SenderAddr() string
	SetBalance(address string, amount ankrcmm.Amount)
	Balance(address string, symbol string) (*big.Int, error)
	SetAllowance(addrSender string, addrSpender string, amount ankrcmm.Amount)
	Allowance(addrSender string, addrSpender string, symbol string) (*big.Int, error)
	BuildCurrencyCAddrMap(symbol string, cAddr string) error
	LoadContract(cAddr string) (*ankrcmm.ContractInfo, error)
	Height() int64
	Publish(ctx context.Context, msg interface{}) error
	PublishWithTags(ctx context.Context, msg interface{}, tags map[string]string) error
}

func SetBCContext(context ContextAKVM) {
	bcContext = context
}

func GetBCContext() ContextAKVM {
	return bcContext
}

type ContextAKVMImpl struct {
	ContextContract
	appstore.ContractStore
	appstore.BCStore
}

func CreateContextAKVM(context ContextContract, appStore appstore.AppStore) ContextAKVM {
	contAKVM :=  &ContextAKVMImpl{context,appStore, appStore }
	bcContext = contAKVM

	return contAKVM
}

