package context

import (
	"github.com/Ankr-network/ankr-chain/account"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	"github.com/tendermint/tendermint/libs/log"
)

type ContextTx interface {
	MinGasPrice() account.Amount
	AppStore() appstore.AppStore
	Logger() log.Logger
}