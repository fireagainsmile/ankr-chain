package contract

import (
	"math/big"

	ankrcmm "github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/context"
	"github.com/Ankr-network/ankr-chain/store/appstore"
)

//ERC20 standard interface
type ContractERC20 interface {
	Name() string
	Symbol() string
	Decimals() uint8
	TotalSupply() *big.Int
	BalanceOf() *big.Int
	Transfer(toAddr string, value *big.Int) bool
	TransferFrom(fromAddr string, toAddr string, value *big.Int) bool
	Approve(spenderAddr string, value *big.Int) bool
	Allowance(ownerAddr string, spenderAddr string)*big.Int
}

type Invoker interface {
	Invoke(context context.ContextContract, conPatt ankrcmm.ContractPatternType, appStore appstore.AppStore, code []byte, contractName string, method string, param []*ankrcmm.Param, rtnType string) (*ankrcmm.ContractResult, error)
}

