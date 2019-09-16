package account

import (
	"encoding/base64"
	"fmt"
	"strconv"


	"github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/common/code"
	ankrcrypto "github.com/Ankr-network/ankr-chain/crypto"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	"github.com/Ankr-network/ankr-chain/tx"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

func NewAdminAccountMsgTxMsg() *tx.TxMsg {
	return  &tx.TxMsg{ImplTxMsg: new(AdminAccountMsg)}
}

type AdminAccountMsg struct {
}

func (s *AdminAccountMsg) GasWanted() int64 {
	return 1
}

func (s *AdminAccountMsg) GasUsed() int64 {
	return 0
}

func (s *AdminAccountMsg) Type() string {
	return ankrtypes.SetOpPrefix
}

func (s *AdminAccountMsg) Bytes() []byte {
	return nil
}
func (s *AdminAccountMsg) SetSecretKey(sk ankrcrypto.SecretKey) {

}

func (s *AdminAccountMsg) SecretKey() ankrcrypto.SecretKey {
	return nil
}

func (s *AdminAccountMsg) ProcessTx(txMsg interface{}, appStore appstore.AppStore, isOnlyCheck bool) (uint32, string, []cmn.KVPair) {
	trxSetOpSlices, ok := txMsg.([]string)
	if !ok {
		return  code.CodeTypeEncodingError, fmt.Sprintf("invalid tx set op msg"), nil
	}

	if len(trxSetOpSlices) != 5{
		return code.CodeTypeEncodingError, fmt.Sprintf("Set Op incorrect format, got %d", len(trxSetOpSlices)), nil
	}

	keynameS  := trxSetOpSlices[0]
	valueS    := trxSetOpSlices[1]
	nonceS    := trxSetOpSlices[2]
	adminPubS := trxSetOpSlices[3]
	sigS      := trxSetOpSlices[4]

	if adminPubS != adminPubKey() {
		return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected pubkey. Got %v", adminPubS), nil
	}

	if keynameS != ankrtypes.ADMIN_OP_VAL_PUBKEY_NAME && keynameS != ankrtypes.ADMIN_OP_FUND_PUBKEY_NAME &&
		keynameS != ankrtypes.ADMIN_OP_METERING_PUBKEY_NAME {
		return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected keyname. Got %v", keynameS), nil
	}

	nonceInt, err_n := strconv.ParseInt(string(nonceS), 10, 64)
	if err_n != nil {
		return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected nonce. Got %v", nonceS), nil
	}

	pDec, _ := base64.StdEncoding.DecodeString(sigS)
	pubKeyObject, err_d := common.DeserilizePubKey(adminPubS)
	if err_d != nil {
		return  code.CodeTypeEncodingError, fmt.Sprintf("Deserilize pubkey failure. Got %v", adminPubS), nil
	}

	s256 := common.ConvertBySha256([]byte(fmt.Sprintf("%s%s%s", keynameS, valueS, nonceS)))
	bb := pubKeyObject.VerifyBytes(s256[:32], pDec)
	if !bb {
		fmt.Println("Bad signature, transaction failed.", sigS)
		return code.CodeTypeEncodingError, fmt.Sprintf("Bad signature. Got %v", sigS), nil
	}

	var inNonce string = "0"
	inNonceByte := appStore.Get([]byte(ankrtypes.SET_OP_NONCE))
	if len(inNonceByte) != 0 {
		inNonce = string(inNonceByte)
	}

	inNonceInt, err_p:= strconv.ParseInt(string(inNonce), 10, 64)
	if err_p != nil || inNonceInt < 0 {
		return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected from nonce. Got %v", inNonce), nil
	}

	if (inNonceInt + 1) != nonceInt {
		return code.CodeTypeEncodingError, fmt.Sprintf("nonce should be one more than last nonce. Got %v", nonceS), nil
	}

	if isOnlyCheck {
		return  code.CodeTypeOK, "", nil
	}

	appStore.Set([]byte(keynameS), []byte(valueS))
	appStore.Set([]byte(ankrtypes.SET_OP_NONCE), []byte(nonceS))

	tags := []cmn.KVPair{
		{Key: []byte("app.type"), Value: []byte("SetOp")},
	}

	return code.CodeTypeOK, "", tags
}

