package decoder

type TxDecoderAdapter struct {
}

func (txda * TxDecoderAdapter) Decode(txBytes []byte) (txType string, data interface{}, err error){
	txType, data, err = new(TxDecoderV0).Decode(txBytes)
	if err == nil {
		return
	}

	txType, data, err = new(TxDecoderV0).Decode(txBytes)

	return
}
