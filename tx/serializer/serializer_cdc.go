package serializer

import (
	"errors"

	"github.com/Ankr-network/ankr-chain/tx"
	"github.com/tendermint/go-amino"
)

type TxSerializerCDC struct {
	txCDC *amino.Codec
}

func NewTxSerializerCDC() *TxSerializerCDC {
	return &TxSerializerCDC {CreateTxCDC()}
}

func (txs *TxSerializerCDC) Serialize(txMsg *tx.TxMsg) ([]byte, error) {
	return txs.txCDC.MarshalBinaryLengthPrefixed(txMsg)
}

func (txs *TxSerializerCDC) MarshalJSON(msg interface{}) ([]byte, error) {
	return txs.txCDC.MarshalJSON(msg)
}

func (txs *TxSerializerCDC) Deserialize(txBytes []byte) (*tx.TxMsg, error) {
	var txMsg tx.TxMsg

	if len(txBytes) == 0 {
		return nil, errors.New("nil tx")
	}

	err := txs.txCDC.UnmarshalBinaryLengthPrefixed(txBytes, &txMsg)
	if err == nil {
		return &txMsg, nil
	} else {
		return nil, err
	}
}