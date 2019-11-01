package common

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"

	"github.com/tendermint/tendermint/crypto/ed25519"
)

const (
	PubKeyEd25519Size = 32
)

func ConvertBySha256(input []byte) [32]byte{
	sum := sha256.Sum256([]byte(input))
	return sum
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

func AddressByPublicKey(pub_key string) (string, error) {
	pubKeyObject, err := DeserilizePubKey(pub_key)
	if err != nil {
		return "", err
	}

	address := fmt.Sprintf("%s", pubKeyObject.Address())
	return address, nil
}

/*
use pem format string ECDSA public_key to verify the input's signature.
*/
func EcdsaVerify(pubpem, input, signature1, signature2 string) bool{

	rSig := new(big.Int)
	rSig, ok := rSig.SetString(signature1, 10)
	if !ok {
		fmt.Println("SetString: error")
		return false
	}

	sSig := new(big.Int)
	sSig, ok = sSig.SetString(signature2, 10)
	if !ok {
		fmt.Println("SetString: error")
		return false
	}

	ecPublicKey, err := parseEcdsaPublicKeyFromPemStr(pubpem)
	if (err != nil) {
		fmt.Println(err)
		return false
	}

	sum := sha256.Sum256([]byte(input))
	valid := ecdsa.Verify(ecPublicKey, sum[:32], rSig, sSig)

	return valid
}

func parseEcdsaPublicKeyFromPemStr(pubPEM string) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		fmt.Println("failed to parse certificate PEM")
		return nil, errors.New("failed to parse PEM block containing the cert")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		fmt.Println("failed to parse certificate: ", err.Error())
		return nil, err
	}

	pub := cert.PublicKey.(*ecdsa.PublicKey)

	return pub, nil
}