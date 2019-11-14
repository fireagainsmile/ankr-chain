package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/tendermint/tendermint/crypto"
	"math/big"

	"github.com/Ankr-network/ankr-chain/common"
)

type SecretKeyPem struct {
	PrivPEM   string
	PEMBase64 string
	Extra     string
}

func ParseEcdsaPrivateKeyFromPemStr(privPEM string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the privkey")
	}

	priv, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

func ParseEcdsaPublicKeyFromPemStr(pubPEM string) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the cert")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	pub := cert.PublicKey.(*ecdsa.PublicKey)

	return pub, nil
}

func NewSecretKeyPem(PrivPEM string, PEMBase64 string, Extra string) *SecretKeyPem {
	return &SecretKeyPem {PrivPEM, PEMBase64, Extra}
}

func (skp *SecretKeyPem) PubKey() (string, error) {
	return "", errors.New("SecretKeyPem not support method PubKey")
}

func (skp *SecretKeyPem) PriKey() (string, error) {
	return "", errors.New("SecretKeyPem not support method PriKey")
}

func (skp *SecretKeyPem) Address() (common.Address, error) {
	pubKey, err := ParseEcdsaPublicKeyFromPemStr(skp.PEMBase64)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := append(pubKey.X.Bytes(), pubKey.Y.Bytes()...)
	if skp.Extra != "" {
		pubKeyBytes = append(pubKeyBytes,skp.Extra...)
	}

	addrHex := crypto.AddressHash(pubKeyBytes)

	return  common.Address(addrHex.String()), err
}

func (skp *SecretKeyPem) Sign(msg []byte) (*Signature, error) {
	privKey, err := ParseEcdsaPrivateKeyFromPemStr(skp.PrivPEM)
	if err != nil {
		return nil, err
	}

	sum := sha256.Sum256(msg)
	r, s, err := ecdsa.Sign(rand.Reader, privKey, sum[:32])
    if err != nil {
    	return nil, err
	}

	return &Signature{R:r.String(), S: s.String(), PubPEM: skp.PEMBase64}, nil
}

func (skp *SecretKeyPem) Verify(msg []byte, signature *Signature) bool {
	if msg == nil || signature == nil {
		return false
	}

	pubKey, err := ParseEcdsaPublicKeyFromPemStr(signature.PubPEM)
	if err != nil {
		return false
	}

	r, isSucess := new(big.Int).SetString(signature.R, 10)
	if !isSucess {
		return false
	}

	s, isSucess := new(big.Int).SetString(signature.S, 10)
	if !isSucess {
		return false
	}

	sum := sha256.Sum256(msg)
	return ecdsa.Verify(pubKey, sum[:32], r, s)
}
