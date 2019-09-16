package account

import (
	acm "github.com/Ankr-network/ankr-chain/common"
)

func adminPubKey() string {
	if acm.RM == acm.RunModeTesting {
		return "sxhP4F6OLZKNPQ2lG13WcHzitNX9++h56cppBDhwMlI="
	}else if acm.RM == acm.RunModeProd {
		return "j90knB4tx3d6xi9KefyCl2FwS/hd/jpEj+cbHdzFcqM="
	}

	return ""
}