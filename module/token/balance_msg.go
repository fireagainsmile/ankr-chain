package token

import (
	"encoding/base64"
	"fmt"
	"github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/common/code"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	"github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"math/big"
	"strconv"
	"strings"
)

type BalanceMsg struct {

}

func (v *BalanceMsg) CheckTx(tx []byte, appStore appstore.AppStore) types.ResponseCheckTx {
	tx = tx[len(ankrtypes.SetBalancePrefix):]
	trxSetBalanceSlices := strings.Split(string(tx), ":")
	if len(trxSetBalanceSlices) != 5{
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Set Balance incorrect format, got %d", len(tx))}
	}

	addressS := trxSetBalanceSlices[0]
	amountS := trxSetBalanceSlices[1]
	nonceS := trxSetBalanceSlices[2]
	adminPubS := trxSetBalanceSlices[3]
	sigS := trxSetBalanceSlices[4]


	var admin_pubkey_str string = ""
	admin_pubkey := appStore.Get([]byte(ankrtypes.ADMIN_OP_FUND_PUBKEY_NAME))
	if len(admin_pubkey) == 0 {
		fmt.Println("use default ADMIN_OP_FUND_PUBKEY_NAME")
		admin_pubkey_str = adminPubkeyOfBalance()
	} else {
		admin_pubkey_str = string(admin_pubkey)
	}

	if adminPubS != admin_pubkey_str {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected pubkey. Got %v", adminPubS)}
	}

	if len(addressS) != ankrtypes.KeyAddressLen {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected address. Got %v", addressS)}
	}

	amountSet, err := new(big.Int).SetString(string(amountS), 10)
	if !err {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected amount. Got %v", amountS)}
	} else { // amountSet < 0
		zeroN, _ := new(big.Int).SetString("0", 10)
		if amountSet.Cmp(zeroN) == -1 {
			return types.ResponseCheckTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Unexpected amount, negative num. Got %v", amountS)}
		}
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

	s256 := common.ConvertBySha256([]byte(fmt.Sprintf("%s%s%s", addressS, amountS, nonceS)))
	bb := pubKeyObject.VerifyBytes(s256[:32], pDec)
	if !bb {
		fmt.Println("Bad signature, transaction failed.", sigS)
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Bad signature. Got %v", sigS)}
	}

	inBalanceAndNonce := appStore.Get(ankrtypes.PrefixBalanceKey([]byte(addressS)))
	balanceNonceSlices := strings.Split(string(inBalanceAndNonce), ":")
	var inBalance string
	var inNonce string
	if len(balanceNonceSlices) == 1 {
		inBalance = balanceNonceSlices[0]
		inNonce = "0"
	} else if len(balanceNonceSlices) == 2 {
		inBalance = balanceNonceSlices[0]
		inNonce = balanceNonceSlices[1]
	} else {
		inBalance = "0"
		inNonce = "0"
	}
	_ = inBalance

	inNonceInt, err_p:= strconv.ParseInt(string(inNonce), 10, 64)
	if err_p != nil || inNonceInt < 0 {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected from nonce. Got %v", inNonce)}
	}

	if (len(balanceNonceSlices) == 2) && (inNonceInt + 1) != nonceInt {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("nonce should be one more than last nonce. Got %v", nonceS)}
	}

	return types.ResponseCheckTx{Code: code.CodeTypeOK, GasWanted: 1}
}

/* will add signature verification when wallet code is ready */
func (v *BalanceMsg) DeliverTx(tx []byte, appStore appstore.AppStore) types.ResponseDeliverTx {
	tx = tx[len(ankrtypes.SetBalancePrefix):]
	trxSetBalanceSlices := strings.Split(string(tx), ":")
	if len(trxSetBalanceSlices) != 5 {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Expected trx set balance. Got %v", trxSetBalanceSlices)}
	}
	addressS := trxSetBalanceSlices[0]
	amountS := trxSetBalanceSlices[1]
	nonceS := trxSetBalanceSlices[2]
	adminPubS := trxSetBalanceSlices[3]
	sigS := trxSetBalanceSlices[4]

	var admin_pubkey_str string = ""
	admin_pubkey := appStore.Get([]byte(ankrtypes.ADMIN_OP_FUND_PUBKEY_NAME))
	if len(admin_pubkey) == 0 {
		fmt.Println("use default ADMIN_OP_FUND_PUBKEY_NAME")
		admin_pubkey_str = adminPubkeyOfBalance()
	} else {
		admin_pubkey_str = string(admin_pubkey)
	}

	if adminPubS != admin_pubkey_str {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected pubkey. Got %v", adminPubS)}
	}

	if len(addressS) != ankrtypes.KeyAddressLen {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected address. Got %v", addressS)}
	}

	amountSet, err := new(big.Int).SetString(string(amountS), 10)
	if !err {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected amount. Got %v", amountS)}
	} else { // amountSet < 0
		zeroN, _ := new(big.Int).SetString("0", 10)
		if amountSet.Cmp(zeroN) == -1 {
			return types.ResponseDeliverTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Unexpected amount, negative num. Got %v", amountS)}
		}
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

	s256 := common.ConvertBySha256([]byte(fmt.Sprintf("%s%s%s", addressS, amountS, nonceS)))
	bb := pubKeyObject.VerifyBytes(s256[:32], pDec)
	if !bb {
		fmt.Println("Bad signature, transaction failed.", sigS)
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Bad signature. Got %v", sigS)}
	}

	inBalanceAndNonce := appStore.Get(ankrtypes.PrefixBalanceKey([]byte(addressS)))
	balanceNonceSlices := strings.Split(string(inBalanceAndNonce), ":")
	var inBalance string
	var inNonce string
	if len(balanceNonceSlices) == 1 {
		inBalance = balanceNonceSlices[0]
		inNonce = "0"
	} else if len(balanceNonceSlices) == 2 {
		inBalance = balanceNonceSlices[0]
		inNonce = balanceNonceSlices[1]
	} else {
		inBalance = "0"
		inNonce = "0"
	}
	_ = inBalance

	inNonceInt, err_p:= strconv.ParseInt(string(inNonce), 10, 64)
	if err_p != nil || inNonceInt < 0 {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected from nonce. Got %v", inNonce)}
	}

	if (len(balanceNonceSlices) == 2) && (inNonceInt + 1) != nonceInt {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("nonce should be one more than last nonce. Got %v", nonceS)}
	}

	appStore.Set(ankrtypes.PrefixBalanceKey([]byte(addressS)), []byte(amountS + ":" + nonceS))
	appStore.IncSize()

	tags := []cmn.KVPair{
		{Key: []byte("app.type"), Value: []byte("SetBalance")},
	}

	return types.ResponseDeliverTx{Code: code.CodeTypeOK,  GasUsed: 0, Tags: tags}
}
