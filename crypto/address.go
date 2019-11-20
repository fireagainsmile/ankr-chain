package crypto

import (
	"github.com/Ankr-network/ankr-chain/common"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

type CertAddrType int
const (
	_ CertAddrType = iota
	CertAddrTypeSet
	CertAddrTypeRemove
)

func CreateContractAddress(callerAddr string, nonce uint64) string {
	hasher := tmhash.NewTruncated()
	hasher.Write([]byte(callerAddr))
	hasher.Write(common.UInt64ToBytes(nonce))
	bytesSum :=  hasher.Sum(nil)

	return  crypto.Address(bytesSum).String()
}

func CreateCertAddress(pubBS64 string, dcName string, addrType CertAddrType) string{
	hasher := tmhash.NewTruncated()

	addr, _ := common.AddressByPublicKey(pubBS64)
	hasher.Write([]byte(addr))
	hasher.Write([]byte(dcName))
	if addrType == CertAddrTypeSet {
		hasher.Write([]byte(dcName + ":set_crt"))
	} else if addrType == CertAddrTypeRemove {
		hasher.Write([]byte(dcName + ":remove_crt"))
	}

	bytesSum :=  hasher.Sum(nil)

	return  crypto.Address(bytesSum).String()
}