package token

import (
	"github.com/Ankr-network/ankr-chain/router"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
)

func init() {
	trMsg  := new(TransferMsg)
	balMsg := new(BalanceMsg)
	router.MsgRouterInstance().AddTxMessageHandler(ankrtypes.TrxSendPrefix, trMsg)
	router.MsgRouterInstance().AddTxMessageHandler(ankrtypes.SetBalancePrefix, balMsg)
}
