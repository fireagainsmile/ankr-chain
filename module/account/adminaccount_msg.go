package account

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	
	"github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/common/code"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	"github.com/tendermint/tendermint/abci/types"
)

type AdminAccountMsg struct {

}

func (ap *AdminAccountMsg) CheckTx(tx []byte, appStore appstore.AppStore) types.ResponseCheckTx {
	tx = tx[len(ankrtypes.SetOpPrefix):]
	trxSetOpSlices := strings.Split(string(tx), ":")
	if len(trxSetOpSlices) != 5{
		return types.ResponseCheckTx{
							Code: code.CodeTypeEncodingError,
							Log:  fmt.Sprintf("Set Balance incorrect format, got %d", len(tx))}
	}

	keynameS := trxSetOpSlices[0]
	valueS := trxSetOpSlices[1]
	nonceS := trxSetOpSlices[2]
	adminPubS := trxSetOpSlices[3]
	sigS := trxSetOpSlices[4]

	if adminPubS != adminPubKey() {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected pubkey. Got %v", adminPubS)}
	}

	if keynameS != ankrtypes.ADMIN_OP_VAL_PUBKEY_NAME && keynameS != ankrtypes.ADMIN_OP_FUND_PUBKEY_NAME &&
			  keynameS != ankrtypes.ADMIN_OP_METERING_PUBKEY_NAME {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected keyname. Got %v", keynameS)}
	}

	nonceInt, err_n := strconv.ParseInt(string(nonceS), 10, 64)
	if err_n != nil {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected nonce. Got %v", nonceS)}
	}

	pDec, _ := base64.StdEncoding.DecodeString(sigS)
	pubKeyObject, err_d := common.DeserilizePubKey(adminPubS)
	if err_d != nil {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Deserilize pubkey failure. Got %v", adminPubS)}
	}

	s256 := common.ConvertBySha256([]byte(fmt.Sprintf("%s%s%s", keynameS, valueS, nonceS)))
	bb := pubKeyObject.VerifyBytes(s256[:32], pDec)
	if !bb {
		fmt.Println("Bad signature, transaction failed.", sigS)
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Bad signature. Got %v", sigS)}
		}

	var inNonce string = "0"
	inNonceByte := appStore.Get([]byte(ankrtypes.SET_OP_NONCE))
	if len(inNonceByte) != 0 {
		inNonce = string(inNonceByte)
	}

	inNonceInt, err_p:= strconv.ParseInt(string(inNonce), 10, 64)
	if err_p != nil || inNonceInt < 0 {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected from nonce. Got %v", inNonce)}
	}

	if (inNonceInt + 1) != nonceInt {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("nonce should be one more than last nonce. Got %v", nonceS)}
	}

	return types.ResponseCheckTx{Code: code.CodeTypeOK, GasWanted: 1}
} 

func (ap *AdminAccountMsg) DeliverTx(tx []byte, appStore appstore.AppStore) types.ResponseDeliverTx {
	tx = tx[len(ankrtypes.SetOpPrefix):]
	trxSetOpSlices := strings.Split(string(tx), ":")
	if len(trxSetOpSlices) != 5{
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Set Balance incorrect format, got %d", len(tx))}
	}

	keynameS := trxSetOpSlices[0]
	valueS := trxSetOpSlices[1]
	nonceS := trxSetOpSlices[2]
	adminPubS := trxSetOpSlices[3]
	sigS := trxSetOpSlices[4]

	if adminPubS != adminPubKey() {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected pubkey. Got %v", adminPubS)}
	}

	if keynameS != ankrtypes.ADMIN_OP_VAL_PUBKEY_NAME && keynameS != ankrtypes.ADMIN_OP_FUND_PUBKEY_NAME &&
		keynameS != ankrtypes.ADMIN_OP_METERING_PUBKEY_NAME {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected keyname. Got %v", keynameS)}
	}

	nonceInt, err_n := strconv.ParseInt(string(nonceS), 10, 64)
	if err_n != nil {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected nonce. Got %v", nonceS)}
	}

	pDec, _ := base64.StdEncoding.DecodeString(sigS)
	pubKeyObject, err_d := common.DeserilizePubKey(adminPubS)
	if err_d != nil {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Deserilize pubkey failure. Got %v", adminPubS)}
	}

	s256 := common.ConvertBySha256([]byte(fmt.Sprintf("%s%s%s", keynameS, valueS, nonceS)))
	bb := pubKeyObject.VerifyBytes(s256[:32], pDec)
	if !bb {
		fmt.Println("Bad signature, transaction failed.", sigS)
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Bad signature. Got %v", sigS)}
	}

	var inNonce string = "0"
	inNonceByte := appStore.Get([]byte(ankrtypes.SET_OP_NONCE))
	if len(inNonceByte) != 0 {
		inNonce = string(inNonceByte)
	}

	inNonceInt, err_p:= strconv.ParseInt(string(inNonce), 10, 64)
	if err_p != nil || inNonceInt < 0 {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected from nonce. Got %v", inNonce)}
	}

	if (inNonceInt + 1) != nonceInt {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("nonce should be one more than last nonce. Got %v", nonceS)}
	}

	appStore.Set([]byte(keynameS), []byte(valueS))
	appStore.Set([]byte(ankrtypes.SET_OP_NONCE), []byte(nonceS))
	appStore.IncSize()

	return types.ResponseDeliverTx{Code: code.CodeTypeOK, GasWanted: 1, GasUsed: 0}
}
