package types

import (
	"github.com/Ankr-network/ankr-chain/account"
	"github.com/tendermint/go-amino"
)

const (
	APPName = "AnkrApp"

	KeyAddressLen = 46
	ValidatorSetChangePrefix string = "val:"
	AccountBlancePrefix string = "bal:"
	AccountStakePrefix string = "stk"
	CertPrefix string = "crt:"
	MeteringPrefix string = "mtr:"
	AllAccountsPrefix string = "all_accounts"
	AllCrtsPrefix string = "all_crts"

	SetMeteringPrefix string = "set_mtr="
	TrxSendPrefix string = "trx_send="
	SetBalancePrefix string = "set_bal="
	SetOpPrefix string = "set_op="
	SetStakePrefix string = "set_stk="
	SetCertPrefix string = "set_crt="
	RemoveCertPrefix string = "rmv_crt="

	SET_CRT_NONCE string = "set_crt_nonce"
	RMV_CRT_NONCE string = "rmv_crt_nonce"
	SET_OP_NONCE string = "admin_nonce"
	SET_VAL_NONCE string = "val_nonce"
	ADMIN_OP_VAL_PUBKEY_NAME string = "admin_op_val_pubkey"
	ADMIN_OP_FUND_PUBKEY_NAME string = "admin_op_fund_pubkey"
	ADMIN_OP_METERING_PUBKEY_NAME string = "admin_op_metering_pubkey"
)

func PrefixBalanceKey(key []byte) []byte {
	return append([]byte(AccountBlancePrefix), key...)
}

// CommitID contains the tree version number and its merkle root.
type CommitID struct {
	Version int64
	Hash    []byte
}

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
	StakeAmount  account.Amount `json:"stakeamount"`
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



