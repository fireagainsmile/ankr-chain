package v0

import (
	"encoding/base64"
	"fmt"
	"github.com/Ankr-network/ankr-chain/crypto"
	"github.com/Ankr-network/新建文件夹/ankr-chain/common"
	"strconv"
	"time"

	"github.com/Ankr-network/ankr-chain/common/code"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	"github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

type meteringMsg struct {}

func (m *meteringMsg) ProcessTx(txMsg interface{}, appStore appstore.AppStore) types.ResponseDeliverTx {
	trxSetMeteringSlices, ok := txMsg.([]string)
	if !ok {
		return types.ResponseDeliverTx{ Code:code.CodeTypeEncodingError, Log: fmt.Sprintf("invalid tx set metering msg") }
	}

	if len(trxSetMeteringSlices) != 6 {
		return types.ResponseDeliverTx{ Code:code.CodeTypeEncodingError, Log: fmt.Sprintf("Expected trx set metering. Got %v", trxSetMeteringSlices) }
	}

	dcS    := trxSetMeteringSlices[0]
	nsS    := trxSetMeteringSlices[1]
	sigxS  := trxSetMeteringSlices[2]
	sigyS  := trxSetMeteringSlices[3]
	nonceS := trxSetMeteringSlices[4]
	valueS := trxSetMeteringSlices[5]

	nonceInt, err_nonce := strconv.ParseUint(string(nonceS), 10, 64)
	if err_nonce != nil {
		return types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("Unexpected nonce8. Got %v", nonceS) }
	}

	pemB64 := appStore.CertKey(dcS)
	if len(pemB64) == 0 {
		return types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("can not find cert file of %s", dcS) }
	}

	key := crypto.NewSecretKeyPem("", pemB64,"@mert:"+"dc1_"+"ns1")

	fromAddr, _ := key.Address()

	/* verify nonce */
	nonce, _ := appStore.Nonce(string(fromAddr))

	if nonceInt != nonce {
		return types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("Unexpected cert nonce. Got %v, Expected %v", nonceS, nonce) }
	}

	/* verify sig */
	pemByte, err := base64.StdEncoding.DecodeString(pemB64)
	if err != nil {
		return types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("pem file decoding error. Got %v", string(pemByte)) }
	}
	pem := string(pemByte)

	bResult := common.EcdsaVerify(pem, dcS+nsS+valueS+nonceS, sigxS, sigyS)
	if !bResult {
		return  types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("metering signature wrong. Got %v,%v", sigxS, sigyS) }
	}

	appStore.SetMetering(dcS , nsS, valueS)
	appStore.IncNonce(string(fromAddr))

	tvalue := time.Now().UnixNano()
	tags := []cmn.KVPair{
		{Key: []byte("app.metering"), Value: []byte(dcS + ":" + nsS)},
		{Key: []byte("app.timestamp"), Value: []byte(strconv.FormatInt(tvalue, 10))},
		{Key: []byte("app.type"), Value: []byte("SetMetering")},
	}

	return types.ResponseDeliverTx{Code: code.CodeTypeOK, GasUsed: 0, Tags: tags}
}
