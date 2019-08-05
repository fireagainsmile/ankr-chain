package metering

import (
	"github.com/Ankr-network/ankr-chain/router"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
)

func init() {
	metMsg     := new(MeteringMsg)
	setCertMsg := new(SetCertMsg)
	rmCertMsg  := new(RemoveCertMsg)
	router.MsgRouterInstance().AddTxMessageHandler(ankrtypes.SetMeteringPrefix, metMsg)
	router.MsgRouterInstance().AddTxMessageHandler(ankrtypes.SetCertPrefix, setCertMsg)
	router.MsgRouterInstance().AddTxMessageHandler(ankrtypes.RemoveCertPrefix, rmCertMsg)
}
