package tx

import (
	"github.com/tendermint/go-amino"
)

var TxCdc *amino.Codec

func init() {
	TxCdc = amino.NewCodec()
	TxCdc.RegisterConcrete(TxMsg{}, "ankrchain/TxMsg", nil)
}
