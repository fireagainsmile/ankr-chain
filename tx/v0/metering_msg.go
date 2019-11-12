package v0

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	ankrcmm "github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/common/code"
	"github.com/Ankr-network/ankr-chain/crypto"
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

	pemByte, err := base64.StdEncoding.DecodeString(pemB64)
	if err != nil {
		return types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("pem file decoding error. Got %v", string(pemByte)) }
	}
	pem := string(pemByte)

	dcNS := fmt.Sprintf("%s_%s", dcS, nsS)

	key := crypto.NewSecretKeyPem("", pem,"@mert:"+dcNS)

	fromAddr, _ := key.Address()

	/* verify nonce */
	fmt.Printf("meteringMsg nonceInt: %d, fromAddr=%s\n", nonceInt, string(fromAddr))
	nonce, _ := appStore.Nonce(string(fromAddr))
	fmt.Printf("meteringMsg nonce: %d, fromAddr=%s\n", nonce, string(fromAddr))
	if nonceInt != nonce + 1 {
		return types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("Unexpected cert nonce. Got %v, Expected %v", nonceS, nonce) }
	}

	/* verify sig */
	bResult := ankrcmm.EcdsaVerify(pem, dcS+nsS+valueS+nonceS, sigxS, sigyS)
	if !bResult {
		return  types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("metering signature wrong. Got %v,%v", sigxS, sigyS) }
	}

	appStore.SetMetering(dcS , nsS, valueS)
	appStore.SetNonce(string(fromAddr), nonceInt+1)

	tvalue := time.Now().UnixNano()
	tags := []cmn.KVPair{
		{Key: []byte("app.metering"), Value: []byte(dcS + ":" + nsS)},
		{Key: []byte("app.timestamp"), Value: []byte(strconv.FormatInt(tvalue, 10))},
		{Key: []byte("app.type"), Value: []byte("SetMetering")},
	}

	return types.ResponseDeliverTx{Code: code.CodeTypeOK, GasUsed: 0, Tags: tags}
}
