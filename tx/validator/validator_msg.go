package validator

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	"math/big"
	"strconv"

	"github.com/Ankr-network/ankr-chain/account"
	"github.com/Ankr-network/ankr-chain/common/code"
	ankrcrypto "github.com/Ankr-network/ankr-chain/crypto"
	"github.com/Ankr-network/ankr-chain/tx"
	txcmm "github.com/Ankr-network/ankr-chain/tx/common"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

const (
	MIN_TOKEN_VALIDATOR = "0"
)

type ValidatorMsg struct {
	Action       uint8                          `json:"action"` //1-create; 2-update; 3-remove
	FromAddress  string                         `json:"fromaddress"`
	Name         string                         `json:"name"`
	PubKey       ankrtypes.ValPubKey            `json:"pubkey"`
	StakeAddress string                         `json:"stakeaddress"`
	StakeAmount  account.Amount                 `json:"stakeamount"`
	ValidHeight  uint64                         `json:"validheight"`
	SetFlag      ankrtypes.ValidatorInfoSetFlag `json:"setflag"`
	gasUsed      *big.Int                       `json:"gasused"`
}

type signValidatorMsg struct {
	Action       uint8                          `json:"action"` //1-create; 2-update; 3-remove
	FromAddress  string                         `json:"fromaddress"`
	Name         string                         `json:"name"`
	PubKey       ankrtypes.ValPubKey            `json:"pubkey"`
	StakeAddress string                         `json:"stakeaddress"`
	StakeAmount  account.Amount                 `json:"stakeamount"`
	ValidHeight  uint64                         `json:"validheight"`
	SetFlag      ankrtypes.ValidatorInfoSetFlag `json:"setflag"`
}

func (sv signValidatorMsg) bytes(txSerializer tx.TxSerializer) ([]byte, error){
	return txSerializer.MarshalJSON(&sv)
}

func (v *ValidatorMsg) SignerAddr() []string {
	return []string {v.FromAddress}
}

func (v *ValidatorMsg) GasWanted() int64 {
	return 0
}

func (v *ValidatorMsg) GasUsed() int64 {
	gasUsed, _ := strconv.ParseInt(MIN_TOKEN_VALIDATOR, 0, 64)

	return gasUsed
}

func (v *ValidatorMsg) SpendGas(gas *big.Int) bool {
	v.gasUsed.Add(v.gasUsed, gas)
	return true
}

func (v *ValidatorMsg) Type() string {
	return txcmm.TxMsgTypeValidator
}

func (v *ValidatorMsg) signMsg() *signValidatorMsg {
	return &signValidatorMsg{v.Action, v.FromAddress, v.Name, v.PubKey, v.StakeAddress, v.StakeAmount, v.ValidHeight, v.SetFlag}
}

func (v *ValidatorMsg) Bytes(txSerializer tx.TxSerializer) []byte {
	bytesSign, _ := v.signMsg().bytes(txSerializer)

	return bytesSign
}
func (v *ValidatorMsg) SetSecretKey(sk ankrcrypto.SecretKey) {

}

func (v *ValidatorMsg) SecretKey() ankrcrypto.SecretKey {
	return &ankrcrypto.SecretKeyEd25519{}
}

func (v *ValidatorMsg) PermitKey(store appstore.AppStore, pubKey []byte) bool {
	adminPubkey := store.Get([]byte(ankrtypes.ADMIN_OP_VAL_PUBKEY_NAME))
	if len(adminPubkey) == 0 {
		adminPubkey = []byte(account.AccountManagerInstance().AdminOpAccount(account.AccountAdminValidator))
	}

	adminPubKeyStr, err := base64.StdEncoding.DecodeString(string(adminPubkey))
	if err != nil {
		return false
	}

	return  bytes.Equal(pubKey, []byte(adminPubKeyStr))
}

func (v *ValidatorMsg) ProcessTx(context tx.ContextTx, isOnlyCheck bool) (uint32, string,  []cmn.KVPair) {
	if len(v.FromAddress) != ankrtypes.KeyAddressLen {
		return  code.CodeTypeInvalidAddress, fmt.Sprintf("ValidatorMsg ProcessTx, unexpected from address. Got %s, addr len=%d", v.FromAddress, len(v.FromAddress)), nil
	}

	if len(v.StakeAddress) != ankrtypes.KeyAddressLen {
		return code.CodeTypeInvalidAddress, fmt.Sprintf("ValidatorMsg ProcessTx, unexpected stake address. Got %s, addr len=%d", v.StakeAddress, len(v.StakeAddress)), nil
	}

	if v.StakeAmount.Cur.Symbol != "ANKR" {
		return code.CodeTypeInvalidStakeCurrency, fmt.Sprintf("ValidatorMsg ProcessTx, invalid stake currency: currency=%s", v.StakeAmount.Cur.Symbol), nil
	}

	bal, err := context.AppStore().Balance(v.StakeAddress, "ANKR")
	if err != nil {
		return code.CodeTypeLoadBalError, fmt.Sprintf("ValidatorMsg ProcessTx, load balance err: address=%s, err=%v", v.StakeAddress, v), nil
	}

	amountTemp, _ := new(big.Int).SetString(MIN_TOKEN_VALIDATOR, 10)
	amountTemp = amountTemp.Add(amountTemp, new(big.Int).SetBytes(v.StakeAmount.Value))

	if bal.Cmp(amountTemp) == 0 || bal.Cmp(amountTemp) == -1 {
		return code.CodeTypeBalNotEnough, fmt.Sprintf("ValidatorMsg ProcessTx, balance not enough, bal=%s, expected=%s", bal.String(), amountTemp.String()), nil
	}

	if isOnlyCheck {
		return code.CodeTypeOK, "", nil
	}

	bal = bal.Sub(bal, amountTemp)
	context.AppStore().SetBalance(v.StakeAddress, account.Amount{account.Currency{"ANKR", 18}, bal.Bytes()})

	pubKeyHandler, err := ankrcrypto.GetValPubKeyHandler(&v.PubKey)
	if err != nil {
		return code.CodeTypeInvalidValidatorPubKey, fmt.Sprintf("can't find the respond crypto pubkey handler:type=%s", v.PubKey.Type), nil
	}

	valInfo := &ankrtypes.ValidatorInfo {v.Name,
		pubKeyHandler.Address().String(),
		v.PubKey,
		ValidatorManagerInstance().Power(&v.StakeAmount),
	    v.StakeAddress,
	    v.StakeAmount,
	    v.ValidHeight,
		}

	switch v.Action {
	case 1:
		ValidatorManagerInstance().CreateValidator(valInfo, context.AppStore())
	case 2:
		ValidatorManagerInstance().UpdateValidator(valInfo, v.SetFlag, context.AppStore())
	case 3:
		ValidatorManagerInstance().RemoveValidator(valInfo.ValAddress, context.AppStore())
	}

	tags := []cmn.KVPair{
		{Key: []byte("app.type"), Value: []byte(txcmm.TxMsgTypeValidator)},
	}

	return  code.CodeTypeOK, "", tags
}
