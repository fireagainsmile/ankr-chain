package context

import (
	"math/big"

	"github.com/Ankr-network/ankr-chain/account"
)

type ContextContract interface {
	SpendGas(gas *big.Int) bool
	GasWanted() *big.Int
	GasUsed() *big.Int
	SenderAddr() string
	SetBalance(address string, amount account.Assert)
	Balance(address string, symbol string) (*big.Int, error)
	SetAllowance(addrSender string, addrSpender string, amount account.Assert)
	Allowance(addrSender string, addrSpender string, symbol string) (*big.Int, error)
}
