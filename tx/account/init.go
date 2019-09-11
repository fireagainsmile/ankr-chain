package account

import (
	"github.com/Ankr-network/ankr-chain/router"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
)

func init() {
	router.MsgRouterInstance().AddTxMessageHandler(ankrtypes.SetOpPrefix, NewAdminAccountMsgTxMsg())
}