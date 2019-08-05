package metering

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/common/code"
	"github.com/Ankr-network/ankr-chain/router"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	"github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

func init() {
	mMsg := new(MeteringMsg)
	router.MsgRouterInstance().AddTxMessageHandler(ankrtypes.SetMeteringPrefix, mMsg)
}

type MeteringMsg struct {
}

func (m *MeteringMsg) prefixSetMeteringKey(key []byte) []byte {
	return append([]byte(ankrtypes.MeteringPrefix), key...)
}

func (m *MeteringMsg) CheckTx(tx []byte, appStore appstore.AppStore) types.ResponseCheckTx {
	tx = tx[len(ankrtypes.SetMeteringPrefix):]
	trxSetMeteringSlices := strings.SplitN(string(tx), ":", 6)
	if len(trxSetMeteringSlices) != 6 {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Metering incorrect format, got %s", string(tx))}
	}

	dcS := trxSetMeteringSlices[0]
	nsS := trxSetMeteringSlices[1]
	sigxS := trxSetMeteringSlices[2]
	sigyS := trxSetMeteringSlices[3]
	nonceS := trxSetMeteringSlices[4]
	valueS := trxSetMeteringSlices[5]

	nonceInt, err_nonce := strconv.ParseInt(string(nonceS), 10, 64)
	if err_nonce != nil {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected nonce. Got %v", nonceS)}
	}

	/* verify nonce */
	var nonceOld int64 = 0
	meteringRec := appStore.Get(m.prefixSetMeteringKey([]byte(dcS + ":" + nsS)))
	if meteringRec == nil || string(meteringRec) == "" {
		nonceOld = 0
	} else {
		trxSetMeteringSlices := strings.SplitN(string(meteringRec), ":", 4)
		if len(trxSetMeteringSlices) != 4 {
			return types.ResponseCheckTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Expected trx set metering check. Got %v", trxSetMeteringSlices)}
		}

		nonceOld, err_nonce = strconv.ParseInt(string(trxSetMeteringSlices[3]), 10, 64)
		if err_nonce != nil {
			return types.ResponseCheckTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Unexpected nonce. Got %v", nonceS)}
		}
	}

	if nonceOld + 1 != nonceInt {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("nonce should be one more than last nonce. Got %v", nonceS)}
	}

	pemB64Byte := appStore.Get(prefixCertKey([]byte(dcS)))
	if len(pemB64Byte) == 0 {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("can not find cert file of %s", dcS)}
	}

	pemByte, err := base64.StdEncoding.DecodeString(string(pemB64Byte))
	if err != nil {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("pem file decoding error. Got %v", string(pemByte))}
	}
	pem := string(pemByte)

	bResult := common.EcdsaVerify(pem, dcS+nsS+valueS+nonceS, sigxS, sigyS)
	if !bResult {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("metering signature wrong. Got %v,%v", sigxS, sigyS)}
	}

	return types.ResponseCheckTx{Code: code.CodeTypeOK, GasWanted: 1}
}

/* will add signature verification when wallet code is ready */
func (m *MeteringMsg) DeliverTx(tx []byte, appStore appstore.AppStore) types.ResponseDeliverTx {
	tx = tx[len(ankrtypes.SetMeteringPrefix):]
	trxSetMeteringSlices := strings.SplitN(string(tx), ":", 6)
	if len(trxSetMeteringSlices) != 6 {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Expected trx set metering. Got %v", trxSetMeteringSlices)}
	}
	dcS := trxSetMeteringSlices[0]
	nsS := trxSetMeteringSlices[1]
	sigxS := trxSetMeteringSlices[2]
	sigyS := trxSetMeteringSlices[3]
	nonceS := trxSetMeteringSlices[4]
	valueS := trxSetMeteringSlices[5]

	nonceInt, err_nonce := strconv.ParseInt(string(nonceS), 10, 64)
	if err_nonce != nil {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected nonce8. Got %v", nonceS)}
	}

	/* verify nonce */
	var nonceOld int64 = 0
	meteringRec := appStore.Get(m.prefixSetMeteringKey([]byte(dcS + ":" + nsS)))
	if meteringRec == nil || string(meteringRec) == "" {
		nonceOld = 0
	} else {
		trxSetMeteringSlices := strings.SplitN(string(meteringRec), ":", 4)
		if len(trxSetMeteringSlices) != 4 {
			return types.ResponseDeliverTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Expected trx set metering. Got %v", trxSetMeteringSlices)}
		}

		nonceOld, err_nonce = strconv.ParseInt(string(trxSetMeteringSlices[3]), 10, 64)
		if err_nonce != nil {
			return types.ResponseDeliverTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Unexpected nonce9. Got %v", nonceS)}
		}
	}

	if nonceOld + 1 != nonceInt {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("nonce should be one more than last nonce. Got %v", nonceS)}
	}

	/* verify sig */
	pemB64Byte := appStore.Get(prefixCertKey([]byte(dcS)))
	if len(pemB64Byte) == 0 {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("can not find cert file of %s", dcS)}
	}

	pemByte, err := base64.StdEncoding.DecodeString(string(pemB64Byte))
	if err != nil {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("pem file decoding error. Got %v", string(pemByte))}
	}
	pem := string(pemByte)

	bResult := common.EcdsaVerify(pem, dcS+nsS+valueS+nonceS, sigxS, sigyS)
	if !bResult {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("metering signature wrong. Got %v,%v", sigxS, sigyS)}
	}

	appStore.Set(m.prefixSetMeteringKey([]byte(dcS + ":" + nsS)),
		[]byte(valueS + ":" + sigxS + ":" + sigyS + ":" + nonceS))
	//fmt.Println(string((prefixSetMeteringKey([]byte(dcS + ":" + nsS)))))
	//fmt.Println(string([]byte(valueS + ":" + sigxS + ":" + sigyS + ":" + sigaS + ":" + sigbS + ":" + nonceS)))
	appStore.IncSize()

	tvalue := time.Now().UnixNano()
	tags := []cmn.KVPair{
		{Key: []byte("app.metering"), Value: []byte(dcS + ":" + nsS)},
		{Key: []byte("app.timestamp"), Value: []byte(strconv.FormatInt(tvalue, 10))},
		{Key: []byte("app.type"), Value: []byte("SetMetering")},
	}
	return types.ResponseDeliverTx{Code: code.CodeTypeOK,  GasUsed: 0, Tags: tags}
}
