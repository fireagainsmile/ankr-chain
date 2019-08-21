package tx

import (
	"math/big"
)

type Currency struct {
	Name    string `json:"name"`
	Decimal int64  `json:"decimal"`
}

type Amount struct {
	Cur   Currency `json:"currency"`
	Value big.Int  `json:"value"`
}

type Balance struct {
	Amounts map[string]Amount `json:"amounts"`
}

type TxFee struct {
	Price Amount `json:"price"`
	Gas   int64  `json:"gas"`
}

func (tf TxFee) Bytes() []byte {
	feeBytes, err := TxCdc.MarshalJSON(tf)
	if err != nil {
		panic(err)
	}

	return feeBytes
}