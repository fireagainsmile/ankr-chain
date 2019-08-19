package token

import (
	acm "github.com/Ankr-network/ankr-chain/common"
)

func adminPubkeyOfBalance() string {
	if acm.RM == acm.RunModeTesting {
		return "dBCzB+l/WYxqk+i54a4addy1XhiIK5t0IAZ5OKtegWY="
	}else if acm.RM == acm.RunModeProd {
		return "sasRoTNPFzpJIHkTILaJaBnhcoC78zJk1Jy3s1/xvAE="
	}

	return ""
}