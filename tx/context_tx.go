package tx

import (
	"context"

	ankrcmm "github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/contract"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	"github.com/tendermint/tendermint/libs/log"
)

type Publisher interface {
	Publish(ctx context.Context, msg interface{}) error
	PublishWithTags(ctx context.Context, msg interface{}, tags map[string]string) error
}

type ContextTx interface {
	MinGasPrice() ankrcmm.Amount
	AppStore() appstore.AppStore
	TxSerializer() TxSerializer
	Contract() contract.Contract
	Publisher() Publisher
	Logger() log.Logger
}