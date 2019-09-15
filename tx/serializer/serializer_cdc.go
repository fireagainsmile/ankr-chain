package serializer

import (
	"errors"
	"github.com/Ankr-network/ankr-chain/tx"
)

type TxSerializerCDC struct {
}

func (txs *TxSerializerCDC) Serialize(txMsg *tx.TxMsg) ([]byte, error) {
	return  tx.TxCdc.MarshalBinaryLengthPrefixed(txMsg)
}

func (txs *TxSerializerCDC) Deserialize(txBytes []byte) (*tx.TxMsg, error) {
	var txMsg tx.TxMsg

	if len(txBytes) == 0 {
		return nil, errors.New("nil tx")
	}

	err := tx.TxCdc.UnmarshalBinaryLengthPrefixed(txBytes, &txMsg)
	if err == nil {
		return &txMsg, nil
	} else {
		return nil, err
	}
}