package account

import (
	"github.com/Ankr-network/ankr-chain/router"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
)

func init() {
	adminAccMsg := new(AdminAccountMsg)
	router.MsgRouterInstance().AddTxMessageHandler(ankrtypes.SetOpPrefix, adminAccMsg)
}