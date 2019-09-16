package crypto

import (
	tmcrypto "github.com/tendermint/tendermint/crypto"
)

type Signature struct {
	tmcrypto.PubKey  `json:"pubkey"`
	Signed  []byte   `json:"signed"`
	PubPEM  string   `json:"pubpem"`
	R       string   `json:"r"`
	S       string   `json:"s"`
}