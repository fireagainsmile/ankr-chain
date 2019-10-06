package metering

/*
import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/common/code"
	ankrcrypto"github.com/Ankr-network/ankr-chain/crypto"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	tx "github.com/Ankr-network/ankr-chain/tx"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

func NewMeteringTxMsg() *tx.TxMsg {
	return &tx.TxMsg{ImplTxMsg: new(MeteringMsg)}
}

type MeteringMsg struct {
}

func (m *MeteringMsg) GasWanted() int64 {
	return 0
}

func (m *MeteringMsg) GasUsed() int64 {
	return 0
}

func (m *MeteringMsg) Type() string {
	return ankrtypes.MeteringPrefix
}

func (m *MeteringMsg) Bytes() []byte {
	return nil
}
func (m *MeteringMsg) SetSecretKey(sk ankrcrypto.SecretKey) {

}

func (m *MeteringMsg) SecretKey() ankrcrypto.SecretKey {
	return nil
}

func (m *MeteringMsg) prefixSetMeteringKey(key []byte) []byte {
	return append([]byte(ankrtypes.MeteringPrefix), key...)
}

func (m *MeteringMsg) ProcessTx(txMsg interface{}, appStore appstore.AppStore, isOnlyCheck bool) (uint32, string, []cmn.KVPair) {
	trxSetMeteringSlices, ok := txMsg.([]string)
	if !ok {
		return  code.CodeTypeEncodingError, fmt.Sprintf("invalid tx set metering msg"), nil
	}

	if len(trxSetMeteringSlices) != 6 {
		return code.CodeTypeEncodingError, fmt.Sprintf("Expected trx set metering. Got %v", trxSetMeteringSlices), nil
	}
	dcS    := trxSetMeteringSlices[0]
	nsS    := trxSetMeteringSlices[1]
	sigxS  := trxSetMeteringSlices[2]
	sigyS  := trxSetMeteringSlices[3]
	nonceS := trxSetMeteringSlices[4]
	valueS := trxSetMeteringSlices[5]

	nonceInt, err_nonce := strconv.ParseInt(string(nonceS), 10, 64)
	if err_nonce != nil {
		return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected nonce8. Got %v", nonceS), nil
	}

	//verify nonce
	var nonceOld int64 = 0
	meteringRec := appStore.Get(m.prefixSetMeteringKey([]byte(dcS + ":" + nsS)))
	if meteringRec == nil || string(meteringRec) == "" {
		nonceOld = 0
	} else {
		trxSetMeteringSlices := strings.SplitN(string(meteringRec), ":", 4)
		if len(trxSetMeteringSlices) != 4 {
			return code.CodeTypeEncodingError, fmt.Sprintf("Expected trx set metering. Got %v", trxSetMeteringSlices), nil
		}

		nonceOld, err_nonce = strconv.ParseInt(string(trxSetMeteringSlices[3]), 10, 64)
		if err_nonce != nil {
			return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected nonce9. Got %v", nonceS), nil
		}
	}

	if nonceOld + 1 != nonceInt {
		return  code.CodeTypeEncodingError, fmt.Sprintf("nonce should be one more than last nonce. Got %v", nonceS), nil
	}

	// verify sig
	pemB64Byte := appStore.Get(prefixCertKey([]byte(dcS)))
	if len(pemB64Byte) == 0 {
		return  code.CodeTypeEncodingError, fmt.Sprintf("can not find cert file of %s", dcS), nil
	}

	pemByte, err := base64.StdEncoding.DecodeString(string(pemB64Byte))
	if err != nil {
		return code.CodeTypeEncodingError, fmt.Sprintf("pem file decoding error. Got %v", string(pemByte)), nil
	}
	pem := string(pemByte)

	bResult := common.EcdsaVerify(pem, dcS+nsS+valueS+nonceS, sigxS, sigyS)
	if !bResult {
		return  code.CodeTypeEncodingError, fmt.Sprintf("metering signature wrong. Got %v,%v", sigxS, sigyS), nil
	}

	if isOnlyCheck {
		return code.CodeTypeOK, "", nil
	}

	appStore.Set(m.prefixSetMeteringKey([]byte(dcS + ":" + nsS)),
		[]byte(valueS + ":" + sigxS + ":" + sigyS + ":" + nonceS))

	tvalue := time.Now().UnixNano()
	tags := []cmn.KVPair{
		{Key: []byte("app.metering"), Value: []byte(dcS + ":" + nsS)},
		{Key: []byte("app.timestamp"), Value: []byte(strconv.FormatInt(tvalue, 10))},
		{Key: []byte("app.type"), Value: []byte("SetMetering")},
	}

	return code.CodeTypeOK, "", tags
}

*/
