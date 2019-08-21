package crypto

import (
	"github.com/Ankr-network/ankr-chain/common"
)

type SecretKey interface {
	PubKey() (string, error)
	PriKey() (string, error)
	Address() (common.Address, error)
	Sign(msg []byte) (*Signature, error)
	Verify(msg []byte, signature *Signature) bool
}




