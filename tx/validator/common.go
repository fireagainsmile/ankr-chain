package validator

import (
	"strings"

	ankrcmm "github.com/Ankr-network/ankr-chain/common"
)

func PrefixStakeKey(key []byte) []byte {
	return append([]byte(ankrcmm.AccountStakePrefix), key...)
}

func isValidatorTx(tx []byte) bool {
	return strings.HasPrefix(string(tx), ankrcmm.ValidatorSetChangePrefix)
}
