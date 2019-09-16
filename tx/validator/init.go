package validator

import (
	"github.com/Ankr-network/ankr-chain/router"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
)

func init() {
	router.MsgRouterInstance().AddTxMessageHandler(ankrtypes.ValidatorSetChangePrefix, NewValidatorTxMsg())
	router.MsgRouterInstance().AddTxMessageHandler(ankrtypes.SetStakePrefix, NewStakeTxMsg())
}

