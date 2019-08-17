package account

import (
	"github.com/Ankr-network/ankr-chain/types"
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
   balance   types.Balance
   codes     map[string][]byte
   codeDescs map[string]string
}




