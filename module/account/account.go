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
   accType   AccountType
   address   string
   balances  map[string]*big.Int
   codes     map[string][]byte
   codeDescs map[string]string
}




