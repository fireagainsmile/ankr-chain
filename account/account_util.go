package account

import (
	"github.com/Ankr-network/ankr-chain/common"
	"github.com/tendermint/go-amino"
)

func RegisterCodec(cdc *amino.Codec) {
	cdc.RegisterConcrete(common.AccountInfo{}, "ankr-chain/account", nil)
}

func EncodeAccount(cdc *amino.Codec, accInfo *common.AccountInfo) []byte {
	return cdc.MustMarshalBinaryBare(accInfo)
}

func DecodeAccount(cdc *amino.Codec, bytes []byte) (accInfo common.AccountInfo) {
	cdc.MustUnmarshalBinaryBare(bytes, &accInfo)
	return
}
