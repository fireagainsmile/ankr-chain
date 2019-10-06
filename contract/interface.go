package contract

import (
	"math/big"

	"github.com/Ankr-network/ankr-chain/context"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
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
	Invoke(context context.ContextContract, code []byte, contractName string, method string, param []*ankrtypes.Param, rtnType string) (interface{}, error)
}

