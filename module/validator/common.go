package validator

import (
	acm "github.com/Ankr-network/ankr-chain/common"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
)

func PrefixStakeKey(key []byte) []byte {
	return append([]byte(ankrtypes.AccountStakePrefix), key...)
}

func adminPubKeyOfValidator() string {
	if acm.RM == acm.RunModeTesting {
		return "trwr09Y8sqIdg2H7vhJFsf4aBowBzqkMOjzAGu2ZF6E="
	}else if acm.RM == acm.RunModeProd {
		return "cGSgVIfAsXWbuWImGxJlNzfqruzuGA+4JXv5gfB0FyY="
	}

	return ""
}