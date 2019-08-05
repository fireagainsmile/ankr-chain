package validator

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/common/code"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	"github.com/tendermint/tendermint/abci/types"
	tmCoreTypes "github.com/tendermint/tendermint/types"
)

type ValidatorMsg struct {
}

func (v *ValidatorMsg) isValidatorTx(tx []byte) bool {
	return strings.HasPrefix(string(tx), ankrtypes.ValidatorSetChangePrefix)
}

func (v *ValidatorMsg) CheckTx(tx []byte, appStore appstore.AppStore) types.ResponseCheckTx {
	//val:public_key:power:nonce:admin_pub:sig
	tx = tx[len(ankrtypes.ValidatorSetChangePrefix):]
	pubKeyAndPower := strings.Split(string(tx), ":")
	if len(pubKeyAndPower) != 5 {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Validator incorrect format, got %s", string(tx))}
	}

	pubkeyS, powerS := pubKeyAndPower[0], pubKeyAndPower[1]
	nonceS := pubKeyAndPower[2]
	adminPubS := pubKeyAndPower[3]
	sigS := pubKeyAndPower[4]


	var admin_pubkey_str string = ""
	admin_pubkey := appStore.Get([]byte(ankrtypes.ADMIN_OP_VAL_PUBKEY_NAME))
	if len(admin_pubkey) == 0 {
		fmt.Println("use default ADMIN_OP_VAL_PUBKEY_NAME")
		admin_pubkey_str = adminPubKeyOfValidator()
	} else {
		admin_pubkey_str = string(admin_pubkey)
	}

	if adminPubS != admin_pubkey_str {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected pubkey. Got %v", adminPubS)}
	}

	powerInt, err := strconv.ParseInt(string(powerS), 10, 64)
	if err != nil {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected power. Got %v", powerS)}
	} else { // power < 0
		if powerInt < 0 {
			return types.ResponseCheckTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Unexpected amount, negative num. Got %v", powerS)}
		}
	}

	curValidatorCount := appStore.TotalValidatorPowers(v.isValidatorTx)
	if (curValidatorCount + int64(powerInt)) > tmCoreTypes.MaxTotalVotingPower {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Total powers %d will reach with the power %d", tmCoreTypes.MaxTotalVotingPower, powerInt)}
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

	s256 := common.ConvertBySha256([]byte(fmt.Sprintf("%s%s%s", pubkeyS, powerS, nonceS)))
	bb := pubKeyObject.VerifyBytes(s256[:32], pDec)
	if !bb {
		fmt.Println("Bad signature, transaction failed.", sigS)
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Bad signature. Got %v", sigS)}
	}

	var inNonceInt int64 = 0
	inNonce := appStore.Get(([]byte(ankrtypes.SET_VAL_NONCE)))
	if len(inNonce) == 0 {
		inNonceInt = 0
	} else {
		inNonceIntValue, err_p := strconv.ParseInt(string(inNonce), 10, 64)
		if err_p != nil || inNonceInt < 0 {
			return types.ResponseCheckTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Unexpected from nonce. Got %v", inNonce)}
		}
		inNonceInt = inNonceIntValue
	}

	if (inNonceInt + 1) != nonceInt {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("nonce should be one more than last nonce. Got %v", nonceS)}
	}

	_, err = hex.DecodeString(pubkeyS)
	if err != nil {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Pubkey (%s) is invalid hex", pubkeyS)}
	}

	_, err = strconv.ParseInt(powerS, 10, 64)
	if err != nil {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Power (%s) is not an int", powerS)}
	}

	return types.ResponseCheckTx{Code: code.CodeTypeOK, GasWanted: 1}
}

func (v *ValidatorMsg) DeliverTx(tx []byte, appStore appstore.AppStore) types.ResponseDeliverTx {
	tx = tx[len(ankrtypes.ValidatorSetChangePrefix):]

	//get the pubkey and power
	pubKeyAndPower := strings.Split(string(tx), ":")
	if len(pubKeyAndPower) != 5 {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Expected 'pubkey/power'. Got %v", pubKeyAndPower)}
	}
	pubkeyS, powerS := pubKeyAndPower[0], pubKeyAndPower[1]
	nonceS := pubKeyAndPower[2]
	adminPubS := pubKeyAndPower[3]
	sigS := pubKeyAndPower[4]

	var admin_pubkey_str string = ""
	admin_pubkey := appStore.Get([]byte(ankrtypes.ADMIN_OP_VAL_PUBKEY_NAME))
	if len(admin_pubkey) == 0 {
		fmt.Println("use default ADMIN_OP_VAL_PUBKEY_NAME")
		admin_pubkey_str = adminPubKeyOfValidator()
	} else {
		admin_pubkey_str = string(admin_pubkey)
	}

	if adminPubS != admin_pubkey_str {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected pubkey. Got %v", adminPubS)}
	}

	powerInt, err := strconv.ParseInt(string(powerS), 10, 64)
	if err != nil {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected power. Got %v", powerS)}
	} else { // power < 0
		if powerInt < 0 {
			return types.ResponseDeliverTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Unexpected amount, negative num. Got %v", powerS)}
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

	s256 := common.ConvertBySha256([]byte(fmt.Sprintf("%s%s%s", pubkeyS, powerS, nonceS)))
	bb := pubKeyObject.VerifyBytes(s256[:32], pDec)
	if !bb {
		fmt.Println("Bad signature, transaction failed.", sigS)
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("DeliverTx, Bad signature. Got %v", sigS)}
	}

	var inNonceInt int64 = 0
	inNonce := appStore.Get(([]byte(ankrtypes.SET_VAL_NONCE)))
	if len(inNonce) == 0 {
		inNonceInt = 0
	} else {
		inNonceIntValue, err_p:= strconv.ParseInt(string(inNonce), 10, 64)
		if err_p != nil || inNonceInt < 0 {
			return types.ResponseDeliverTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Unexpected from nonce. Got %v", inNonce)}
		}
		inNonceInt = inNonceIntValue
	}

	if (inNonceInt + 1) != nonceInt {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("nonce should be one more than last nonce. Got %v", nonceS)}
	}

	// decode the pubkey
	pubkey, err := hex.DecodeString(pubkeyS)
	if err != nil {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Pubkey (%s) is invalid hex", pubkeyS)}
	}

	// decode the power
	power, err := strconv.ParseInt(powerS, 10, 64)
	if err != nil {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Power (%s) is not an int", powerS)}
	}

	// update
	appStore.Set([]byte(ankrtypes.SET_VAL_NONCE), []byte(nonceS))
	return ValidatorManagerInstance().UpdateValidator(types.Ed25519ValidatorUpdate(pubkey, int64(power)), appStore)
}