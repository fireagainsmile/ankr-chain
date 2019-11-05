package crypto

import (
	"github.com/Ankr-network/ankr-chain/common"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

func CreateContractAddress(callerAddr string, nonce uint64) string {
	hasher := tmhash.NewTruncated()
	hasher.Write([]byte(callerAddr))
	hasher.Write(common.UInt64ToBytes(nonce))
	bytesSum :=  hasher.Sum(nil)

	return  crypto.Address(bytesSum).String()
}