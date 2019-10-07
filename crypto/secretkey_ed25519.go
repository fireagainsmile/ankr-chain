package crypto

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/Ankr-network/ankr-chain/common"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

const (
	PrivKeyEd25519Size = 64
	PubKeyEd25519Size  = 32
)

func DeserilizePrivKey(priv_key_b64 string) (ed25519.PrivKeyEd25519, error) {
	kDec, err := base64.StdEncoding.DecodeString(priv_key_b64)
	if err != nil {
		return ed25519.PrivKeyEd25519{}, err
	}

	pp := []byte(kDec)
	keyObject := ed25519.PrivKeyEd25519{pp[0], pp[1], pp[2], pp[3], pp[4], pp[5], pp[6], pp[7], pp[8], pp[9],
		pp[10], pp[11], pp[12], pp[13], pp[14], pp[15], pp[16], pp[17], pp[18], pp[19], pp[20], pp[21],
		pp[22], pp[23], pp[24], pp[25], pp[26], pp[27], pp[28], pp[29], pp[30], pp[31], pp[32], pp[33],
		pp[34], pp[35], pp[36], pp[37], pp[38], pp[39], pp[40], pp[41], pp[42], pp[43], pp[44], pp[45],
		pp[46], pp[47], pp[48], pp[49], pp[50], pp[51], pp[52], pp[53], pp[54], pp[55], pp[56], pp[57],
		pp[58], pp[59], pp[60], pp[61], pp[62], pp[PrivKeyEd25519Size-1]}

	return keyObject, nil
}

func DeserilizePubKey(pub_key_b64 string) (ed25519.PubKeyEd25519, error) {
	pDec, err := base64.StdEncoding.DecodeString(pub_key_b64)
	if err != nil {
		return ed25519.PubKeyEd25519{}, err
	}

	pk := []byte(pDec)
	var pubObject ed25519.PubKeyEd25519 = ed25519.PubKeyEd25519{pk[0], pk[1], pk[2], pk[3],pk[4], pk[5],pk[6],
		pk[7],pk[8], pk[9], pk[10], pk[11], pk[12], pk[13], pk[14], pk[15], pk[16], pk[17], pk[18], pk[19],
		pk[20], pk[21],pk[22], pk[23],pk[24], pk[25],pk[26], pk[27],pk[28], pk[29],pk[30], pk[PubKeyEd25519Size - 1]}

	return pubObject, nil
}

type SecretKeyEd25519 struct {
	PrivKey string
}

func NewSecretKeyEd25519(privKey string) *SecretKeyEd25519 {
	return &SecretKeyEd25519{privKey}
}

func (sked *SecretKeyEd25519 ) PubKey() (string, error) {
	privKeyObj, err := DeserilizePrivKey(sked.PrivKey)
	if err != nil {
		return "", err
	}

	return string(privKeyObj.PubKey().Bytes()), nil
}

func (sked *SecretKeyEd25519) PriKey() (string, error) {
	return sked.PrivKey, nil
}

func (sked *SecretKeyEd25519) Address() (common.Address, error) {
	privKeyObj, err := DeserilizePrivKey(sked.PrivKey)
	if err != nil {
		return nil, err
	}

	return common.Address(privKeyObj.PubKey().Address()), nil
}

func (sked *SecretKeyEd25519) Sign(msg []byte) (*Signature, error) {
	privKeyObj, err := DeserilizePrivKey(sked.PrivKey)
	if err != nil {
		return  nil, err
	}

	sum := sha256.Sum256(msg)
	signedBytes, err := privKeyObj.Sign(sum[:32])
	if err != nil {
		return nil, err
	}

	return &Signature{PubKey: privKeyObj.PubKey(), Signed: signedBytes}, nil
}

func (sked *SecretKeyEd25519) Verify(msg []byte, signature *Signature) bool {
	if signature == nil {
		return false
	}

	addr := signature.PubKey.Address()
	if len(addr.String()) != ankrtypes.KeyAddressLen {
		return false
	}

	sum := sha256.Sum256(msg)
	return signature.VerifyBytes(sum[:32], signature.Signed)
}


