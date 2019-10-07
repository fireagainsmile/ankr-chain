package metering

import (
	acm "github.com/Ankr-network/ankr-chain/common"
)

func defaultAdminPubKeyOfMetering() string {
	if acm.RM == acm.RunModeTesting {
		return "wvHG3EddBbXQHcyJal0CS/YQcNYtEbFYxejnqf9OhM4="
	}else if acm.RM == acm.RunModeProd {
		return "cOKct2+weTftBpTvhvFKqzg9tBkN7gG/gtFVuoE53e0="
	}

	return ""
}