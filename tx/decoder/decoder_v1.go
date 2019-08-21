package decoder

import (
	"errors"
	"github.com/Ankr-network/ankr-chain/tx"
)

type TxDecoderV1 struct {
}

func (txdv1 *TxDecoderV1) Decode(txBytes []byte) (txType string, data interface{}, err error) {
	var txMsg tx.TxMsg

	if len(txBytes) == 0 {
		txType = ""
		data   = txMsg
		err    = errors.New("nil tx")
		return
	}

	err = tx.TxCdc.UnmarshalBinaryLengthPrefixed(txBytes, &txMsg)
	if err == nil {
		txType = txMsg.Data.Type()
	} else {
		txType = ""
	}

	return
}