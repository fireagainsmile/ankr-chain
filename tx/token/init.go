package token

import (
	"github.com/Ankr-network/ankr-chain/router"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
)

func init() {
	balMsg := new(BalanceMsg)
	router.MsgRouterInstance().AddTxMessageHandler(ankrtypes.TrxSendPrefix, NewTransferTxM())
	router.MsgRouterInstance().AddTxMessageHandler(ankrtypes.SetBalancePrefix, balMsg)
}
