package handler

import (
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
)

// amino codec to marshal/unmarshal
type Codec = amino.Codec

func New() *Codec {
	cdc := amino.NewCodec()
	return cdc
}

func RegisterCrypto(cdc *Codec) {
	cryptoAmino.RegisterAmino(cdc)
}
var Cdc *Codec

func init() {
	cdc := New()
	RegisterCrypto(cdc)
	Cdc = cdc.Seal()
}
