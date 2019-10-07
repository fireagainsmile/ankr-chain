package validator

import (
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	"strings"
)

func PrefixStakeKey(key []byte) []byte {
	return append([]byte(ankrtypes.AccountStakePrefix), key...)
}

func isValidatorTx(tx []byte) bool {
	return strings.HasPrefix(string(tx), ankrtypes.ValidatorSetChangePrefix)
}
