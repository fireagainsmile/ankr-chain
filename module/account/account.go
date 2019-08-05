package account

import (
	"math/big"
)

type AccountType int
const (
	_ AccountType = iota
	AccountGeneral
	AccountContract
)

type Account struct {
   accountType AccountType
   address     string
   balances    map[string]*big.Int
}


