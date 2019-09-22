package contract

import (
	"github.com/Ankr-network/ankr-chain/contract/native"
	"math/big"

	"github.com/Ankr-network/ankr-chain/context"
	"github.com/Ankr-network/ankr-chain/types"
	"github.com/tendermint/tendermint/libs/log"
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
	Invoke(contractName string, method string,  params []*types.Param) (interface{}, error)
}

func NewInvoker(context context.ContextContract, log log.Logger) Invoker {
	return native.NewNativeInvoker(context, log)
}