package common

import (
	"github.com/tendermint/go-amino"
)

type ValPubKey struct {
	Type string  `json:"type"`
	Data []byte  `json:"data"`
}

type ValidatorInfoSetFlag uint32
const (
	_ ValidatorInfoSetFlag = iota
	ValidatorInfoSetName         = 0x01
	ValidatorInfoSetValAddress   = 0x02
	ValidatorInfoSetPubKey       = 0x04
	ValidatorInfoSetStakeAddress = 0x08
	ValidatorInfoSetStakeAmount  = 0x10
	ValidatorInfoSetValidHeight  = 0x20

)

type ValidatorInfo struct {
	Name         string         `json:"name"`
	ValAddress   string         `json:"valaddress"`
	PubKey       ValPubKey      `json:"pubkey"`
	Power        int64          `json:"power"`
	StakeAddress string         `json:"stakeaddress"`
	StakeAmount  Amount        `json:"stakeamount"`
	ValidHeight  uint64         `json:"validheight"`
}

func EncodeValidatorInfo(cdc *amino.Codec, valInfo *ValidatorInfo) []byte {
	jsonBytes, _ := cdc.MarshalJSON(valInfo)

	return jsonBytes
}

func DecodeValidatorInfo(cdc *amino.Codec, bytes []byte) (valInfo ValidatorInfo) {
	cdc.UnmarshalJSON(bytes, &valInfo)
	return
}