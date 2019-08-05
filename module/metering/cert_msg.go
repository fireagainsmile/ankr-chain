package metering

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

type SetCertMsg struct {

}

func (c *SetCertMsg) CheckTx(tx []byte, appStore appstore.AppStore) types.ResponseCheckTx {
	tx = tx[len(ankrtypes.SetCertPrefix):]
	trxSetCertSlices := strings.SplitN(string(tx), ":", 4)
	if len(trxSetCertSlices) != 4 {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Expected trx set cert. Got %v", trxSetCertSlices)}
	}
	dcS := trxSetCertSlices[0]
	pemB64S := trxSetCertSlices[1]
	nonceS := trxSetCertSlices[2]
	sigS := trxSetCertSlices[3]

	nonceInt, err_nonce := strconv.ParseInt(string(nonceS), 10, 64)
	if err_nonce != nil {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected cert nonce. Got %v, %v", nonceS, err_nonce)}
	}

	nonceOldByte := appStore.Get([]byte(ankrtypes.SET_CRT_NONCE))
	nonceOld, err_nonce := strconv.ParseInt(string(nonceOldByte), 10, 64)
	if err_nonce != nil {
		if len(string(nonceOldByte)) == 0 {
			nonceOld = 0
		} else {
			return types.ResponseCheckTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Unexpected nonce5. Got %v", nonceOld)}
		}
	}

	if nonceOld + 1 != nonceInt {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("nonce should be one more than last nonce. Got %v", nonceS)}
	}

	var admin_pubkey_str string = ""
	admin_pubkey := appStore.Get([]byte(ankrtypes.ADMIN_OP_METERING_PUBKEY_NAME))
	if len(admin_pubkey) == 0 {
		fmt.Println("use default ADMIN_OP_METERING_PUBKEY_NAME")
		admin_pubkey_str = adminPubKeyOfMetering()
	} else {
		admin_pubkey_str = string(admin_pubkey)
	}

	pDec, _ := base64.StdEncoding.DecodeString(sigS)
	pubKeyObject, err := common.DeserilizePubKey(admin_pubkey_str) //set by super user
	if err != nil {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Deserilize pubkey failure. Got %v", admin_pubkey_str)}
	}

	s256 := common.ConvertBySha256([]byte(fmt.Sprintf("%s%s%s", dcS, pemB64S, nonceS)))
	bb := pubKeyObject.VerifyBytes(s256[:32], pDec)
	if !bb {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Bad signature. Got %v", sigS)}
	}
	return types.ResponseCheckTx{Code: code.CodeTypeOK, GasWanted: 1}
}

func (c *SetCertMsg) DeliverTx(tx []byte, appStore appstore.AppStore) types.ResponseDeliverTx {
	tx = tx[len(ankrtypes.SetCertPrefix):]
	trxSetCertSlices := strings.SplitN(string(tx), ":", 4)
	if len(trxSetCertSlices) != 4 {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Expected trx set cert. Got %v", trxSetCertSlices)}
	}
	dcS := trxSetCertSlices[0]
	pemB64S := trxSetCertSlices[1]
	nonceS := trxSetCertSlices[2]
	sigS := trxSetCertSlices[3]

	nonceInt, err_nonce := strconv.ParseInt(string(nonceS), 10, 64)
	if err_nonce != nil {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected cert nonce. Got %v, %v", nonceS, err_nonce)}
	}

	nonceOldByte := appStore.Get([]byte(ankrtypes.SET_CRT_NONCE))
	nonceOld, err_nonce := strconv.ParseInt(string(nonceOldByte), 10, 64)
	if err_nonce != nil {
		if len(string(nonceOldByte)) == 0 {
			nonceOld = 0
		} else {
			return types.ResponseDeliverTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Unexpected nonce5. Got %v", nonceOld)}
		}
	}

	if nonceOld + 1 != nonceInt {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("nonce should be one more than last nonce. Got %v", nonceS)}
	}

	var admin_pubkey_str string = ""
	admin_pubkey := appStore.Get([]byte(ankrtypes.ADMIN_OP_METERING_PUBKEY_NAME))
	if len(admin_pubkey) == 0 {
		fmt.Println("use default ADMIN_OP_METERING_PUBKEY_NAME")
		admin_pubkey_str = adminPubKeyOfMetering()
	} else {
		admin_pubkey_str = string(admin_pubkey)
	}

	pDec, _ := base64.StdEncoding.DecodeString(sigS)
	pubKeyObject, err := common.DeserilizePubKey(admin_pubkey_str) //set by super user
	if err != nil {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Deserilize pubkey failure. Got %v", admin_pubkey_str)}
	}

	s256 := common.ConvertBySha256([]byte(fmt.Sprintf("%s%s%s", dcS, pemB64S, nonceS)))
	bb := pubKeyObject.VerifyBytes(s256[:32], pDec)
	if !bb {
		fmt.Println("Bad signature, transaction failed.", sigS)
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Bad signature. Got %v", sigS)}
	}

	appStore.Set(([]byte(ankrtypes.SET_CRT_NONCE)) ,[]byte(nonceS))
	appStore.Set(prefixCertKey([]byte(dcS)), []byte(pemB64S))
	appStore.IncSize()

	return types.ResponseDeliverTx{Code: code.CodeTypeOK,  GasUsed: 0}
}

type RemoveCertMsg struct {

}

func (c *RemoveCertMsg) CheckTx(tx []byte, appStore appstore.AppStore) types.ResponseCheckTx {
	tx = tx[len(ankrtypes.RemoveCertPrefix):]
	trxSetCertSlices := strings.SplitN(string(tx), ":", 3)
	if len(trxSetCertSlices) != 3 {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Expected trx remove cert. Got %v", trxSetCertSlices)}
	}
	dcS := trxSetCertSlices[0]
	nonceS := trxSetCertSlices[1]
	sigS := trxSetCertSlices[2]

	nonceInt, err_nonce := strconv.ParseInt(string(nonceS), 10, 64)
	if err_nonce != nil {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected nonce6. Got %v", nonceS)}
	}

	nonceOldByte := appStore.Get([]byte(ankrtypes.RMV_CRT_NONCE))
	nonceOld, err_nonce := strconv.ParseInt(string(nonceOldByte), 10, 64)
	if err_nonce != nil {
		if len(string(nonceOldByte)) == 0 {
			nonceOld = 0
		} else {
			return types.ResponseCheckTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Unexpected nonce. Got %v", nonceOld)}
		}
	}

	if nonceOld + 1 != nonceInt {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("nonce should be one more than last nonce. Got %v", nonceS)}
	}

	var admin_pubkey_str string = ""
	admin_pubkey := appStore.Get([]byte(ankrtypes.ADMIN_OP_METERING_PUBKEY_NAME))
	if len(admin_pubkey) == 0 {
		fmt.Println("use default ADMIN_OP_METERING_PUBKEY_NAME")
		admin_pubkey_str = adminPubKeyOfMetering()
	} else {
		admin_pubkey_str = string(admin_pubkey)
	}

	pDec, _ := base64.StdEncoding.DecodeString(sigS)
	pubKeyObject, err := common.DeserilizePubKey(admin_pubkey_str) //set by super user
	if err != nil {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Deserilize pubkey failure. Got %v", admin_pubkey_str)}
	}

	s256 := common.ConvertBySha256([]byte(fmt.Sprintf("%s%s", dcS, nonceS)))
	bb := pubKeyObject.VerifyBytes(s256[:32], pDec)
	if !bb {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Bad signature. Got %v", sigS)}
	}
	return types.ResponseCheckTx{Code: code.CodeTypeOK, GasWanted: 1}
}

func (c *RemoveCertMsg) DeliverTx(tx []byte, appStore appstore.AppStore) types.ResponseDeliverTx {
	tx = tx[len(ankrtypes.RemoveCertPrefix):]
	trxSetCertSlices := strings.SplitN(string(tx), ":", 3)
	if len(trxSetCertSlices) != 3 {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Expected trx remove cert. Got %v", trxSetCertSlices)}
	}
	dcS := trxSetCertSlices[0]
	nonceS := trxSetCertSlices[1]
	sigS := trxSetCertSlices[2]

	nonceInt, err_nonce := strconv.ParseInt(string(nonceS), 10, 64)
	if err_nonce != nil {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected nonce6. Got %v", nonceS)}
	}

	nonceOldByte := appStore.Get([]byte(ankrtypes.RMV_CRT_NONCE))
	nonceOld, err_nonce := strconv.ParseInt(string(nonceOldByte), 10, 64)
	if err_nonce != nil {
		if len(string(nonceOldByte)) == 0 {
			nonceOld = 0
		} else {
			return types.ResponseDeliverTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Unexpected nonce. Got %v", nonceOld)}
		}
	}

	if nonceOld + 1 != nonceInt {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("nonce should be one more than last nonce. Got %v", nonceS)}
	}

	// verify sig
	var admin_pubkey_str string = ""
	admin_pubkey := appStore.Get([]byte(ankrtypes.ADMIN_OP_METERING_PUBKEY_NAME))
	if len(admin_pubkey) == 0 {
		fmt.Println("use default ADMIN_OP_METERING_PUBKEY_NAME")
		admin_pubkey_str = adminPubKeyOfMetering()
	} else {
		admin_pubkey_str = string(admin_pubkey)
	}

	pDec, _ := base64.StdEncoding.DecodeString(sigS)
	pubKeyObject, err := common.DeserilizePubKey(admin_pubkey_str) //set by super user
	if err != nil {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Deserilize pubkey failure. Got %v", admin_pubkey_str)}
	}

	s256 := common.ConvertBySha256([]byte(fmt.Sprintf("%s%s", dcS, nonceS)))
	bb := pubKeyObject.VerifyBytes(s256[:32], pDec)
	if !bb {
		fmt.Println("Bad signature, transaction failed.", sigS)
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Bad signature. Got %v", sigS)}
	}

	appStore.Set(([]byte(ankrtypes.RMV_CRT_NONCE)), []byte(nonceS))
	appStore.Delete(prefixCertKey([]byte(dcS)))
	appStore.IncSize()

	return types.ResponseDeliverTx{Code: code.CodeTypeOK, GasUsed: 0}
}

