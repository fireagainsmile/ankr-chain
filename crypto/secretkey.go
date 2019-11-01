package crypto

import (
	"fmt"

	ankrcmm "github.com/Ankr-network/ankr-chain/common"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type SecretKey interface {
	PubKey() (string, error)
	PriKey() (string, error)
	Address() (ankrcmm.Address, error)
	Sign(msg []byte) (*Signature, error)
	Verify(msg []byte, signature *Signature) bool
}

func GetValPubKeyHandler(valPubkey *ankrcmm.ValPubKey) (tmcrypto.PubKey, error) {
	switch valPubkey.Type {
	case CryptoED25519:
		if len(valPubkey.Data) != ed25519.PubKeyEd25519Size {
			return new(ed25519.PubKeyEd25519), fmt.Errorf("invalid valPubkey data size: type=%s, %d", valPubkey.Type, len(valPubkey.Data))
		}
		var key ed25519.PubKeyEd25519
		copy(key[:], valPubkey.Data)
		return key, nil
	case CryptoSECP256K1:
		if len(valPubkey.Data) != secp256k1.PubKeySecp256k1Size {
			return new(secp256k1.PubKeySecp256k1),  fmt.Errorf("invalid valPubkey data size: type=%s, %d", valPubkey.Type, len(valPubkey.Data))
		}
		var key secp256k1.PubKeySecp256k1
		copy(key[:], valPubkey.Data)
		return key, nil
	default:
		return nil, fmt.Errorf("invalid crypto type: %s", valPubkey.Type)
	}
}




