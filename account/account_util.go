package account

import (
	"github.com/tendermint/go-amino"
)

func RegisterCodec(cdc *amino.Codec) {
	cdc.RegisterConcrete(AccountInfo{}, "ankr-chain/account", nil)
}

func EncodeAccount(cdc *amino.Codec, accInfo *AccountInfo) []byte {
	return cdc.MustMarshalBinaryBare(accInfo)
}

func DecodeAccount(cdc *amino.Codec, bytes []byte) (accInfo AccountInfo) {
	cdc.MustUnmarshalBinaryBare(bytes, &accInfo)
	return
}
