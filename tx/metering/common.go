package metering

import (
	acm "github.com/Ankr-network/ankr-chain/common"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
)

func prefixCertKey(key []byte) []byte {
	return append([]byte(ankrtypes.CertPrefix), key...)
}

func adminPubKeyOfMetering() string {
	if acm.RM == acm.RunModeTesting {
		return "wvHG3EddBbXQHcyJal0CS/YQcNYtEbFYxejnqf9OhM4="
	}else if acm.RM == acm.RunModeProd {
		return "cOKct2+weTftBpTvhvFKqzg9tBkN7gG/gtFVuoE53e0="
	}

	return ""
}