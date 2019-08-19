package validator

import (
	"github.com/Ankr-network/ankr-chain/router"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
)

func init() {
	valMsg  := new(ValidatorMsg)
	stMsg := new(StakeMsg)
	router.MsgRouterInstance().AddTxMessageHandler(ankrtypes.ValidatorSetChangePrefix, valMsg)
	router.MsgRouterInstance().AddTxMessageHandler(ankrtypes.SetStakePrefix, stMsg)
}

