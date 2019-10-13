package tx

import (
	"github.com/Ankr-network/ankr-chain/account"
	"github.com/Ankr-network/ankr-chain/contract"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	"github.com/tendermint/tendermint/libs/log"
	tmpubsub "github.com/tendermint/tendermint/libs/pubsub"
)

type ContextTx interface {
	MinGasPrice() account.Amount
	AppStore() appstore.AppStore
	TxSerializer() TxSerializer
	Contract() contract.Contract
	PubSubServer() *tmpubsub.Server
	Logger() log.Logger
}