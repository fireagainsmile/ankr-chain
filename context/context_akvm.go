package context

import (
	"math/big"

	"github.com/Ankr-network/ankr-chain/account"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
)

var bcContext ContextAKVM

type ContextAKVM interface {
	SpendGas(gas *big.Int) bool
	SenderAddr() string
	SetBalance(address string, amount account.Amount)
	Balance(address string, symbol string) (*big.Int, error)
	SetAllowance(addrSender string, addrSpender string, amount account.Amount)
	Allowance(addrSender string, addrSpender string, symbol string) (*big.Int, error)
	BuildCurrencyCAddrMap(symbol string, cAddr string)
	LoadContract(cAddr string) (*ankrtypes.ContractInfo, error)
	Height() int64
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

