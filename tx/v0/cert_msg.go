package v0

import (
	"encoding/base64"
	"fmt"
	"github.com/Ankr-network/ankr-chain/crypto"
	"strconv"

	"github.com/Ankr-network/ankr-chain/account"
	ankrcmm "github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/common/code"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	txcmm "github.com/Ankr-network/ankr-chain/tx/common"
	"github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

type setCertMsg struct {}

func (sc *setCertMsg) ProcessTx(txMsg interface{}, appStore appstore.AppStore) types.ResponseDeliverTx {
	trxSetCertSlices, ok := txMsg.([]string)
	if !ok {
		return types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("invalid tx set cert msg") }
	}

	if len(trxSetCertSlices) != 4 {
		return types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("Expected trx set cert. Got %v", trxSetCertSlices) }
	}
	dcS     := trxSetCertSlices[0]
	pemB64S := trxSetCertSlices[1]
	nonceS  := trxSetCertSlices[2]
	sigS    := trxSetCertSlices[3]

	nonceInt, err_nonce := strconv.ParseUint(string(nonceS), 10, 64)
	if err_nonce != nil {
		return types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("Unexpected cert nonce. Got %v, %v", nonceS, err_nonce) }
	}

	var admin_pubkey_str = ""
	admin_pubkey := appStore.Get([]byte(ankrcmm.ADMIN_OP_METERING_PUBKEY_NAME))
	if len(admin_pubkey) == 0 {
		//use default ADMIN_OP_METERING_PUBKEY_NAME
		admin_pubkey_str = account.AccountManagerInstance().AdminOpAccount(ankrcmm.AccountAdminMetering)
	} else {
		admin_pubkey_str = string(admin_pubkey)
	}

	addrFrom := crypto.CreateCertAddress(admin_pubkey_str,"dc1")
	nonce, _ := appStore.Nonce(addrFrom)

	if nonceInt != nonce + 1 {
		return types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("Unexpected cert nonce. Got %v, Expected %v", nonceS, nonce) }
	}

	pDec, _ := base64.StdEncoding.DecodeString(sigS)
	pubKeyObject, err := ankrcmm.DeserilizePubKey(admin_pubkey_str) //set by super user
	if err != nil {
		return  types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("Deserilize pubkey failure. Got %v", admin_pubkey_str) }
	}

	s256 := ankrcmm.ConvertBySha256([]byte(fmt.Sprintf("%s%s%s", dcS, pemB64S, nonceS)))
	bb := pubKeyObject.VerifyBytes(s256[:32], pDec)
	if !bb {
		fmt.Println("Bad signature, transaction failed.", sigS)
		return types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("Bad signature. Got %v", sigS) }
	}

	appStore.SetCertKey(dcS, pemB64S)

	appStore.SetNonce(addrFrom, nonce+1)

	tags := []cmn.KVPair{
		{Key: []byte("app.type"), Value: []byte(txcmm.TxMsgTypeSetCertMsg)},
	}

	return types.ResponseDeliverTx{Code: code.CodeTypeOK, GasUsed: 0, Tags: tags}
}

type removeCertMsg struct { }

func (rc *removeCertMsg) ProcessTx(txMsg interface{}, appStore appstore.AppStore) types.ResponseDeliverTx {
	trxSetCertSlices := txMsg.([]string)
	if len(trxSetCertSlices) != 3 {
		return types.ResponseDeliverTx { Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("Expected trx remove cert. Got %v", trxSetCertSlices) }
	}
	dcS    := trxSetCertSlices[0]
	nonceS := trxSetCertSlices[1]
	sigS   := trxSetCertSlices[2]

	nonceInt, err_nonce := strconv.ParseUint(string(nonceS), 10, 64)
	if err_nonce != nil {
		return types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("Unexpected nonce6. Got %v", nonceS) }
	}

	var admin_pubkey_str = ""
	admin_pubkey := appStore.Get([]byte(ankrcmm.ADMIN_OP_METERING_PUBKEY_NAME))
	if len(admin_pubkey) == 0 {
		//use default ADMIN_OP_METERING_PUBKEY_NAME
		admin_pubkey_str = account.AccountManagerInstance().AdminOpAccount(ankrcmm.AccountAdminMetering)
	} else {
		admin_pubkey_str = string(admin_pubkey)
	}

	addrFrom := crypto.CreateCertAddress(admin_pubkey_str,"dc1")
	nonce, _ := appStore.Nonce(addrFrom)

	if nonceInt != nonce + 1 {
		return types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("Unexpected cert nonce. Got %v, Expected %v", nonceS, nonce) }
	}

	pDec, _ := base64.StdEncoding.DecodeString(sigS)
	pubKeyObject, err := ankrcmm.DeserilizePubKey(admin_pubkey_str) //set by super user
	if err != nil {
		return types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("Deserilize pubkey failure. Got %v", admin_pubkey_str) }
	}

	s256 := ankrcmm.ConvertBySha256([]byte(fmt.Sprintf("%s%s", dcS, nonceS)))
	bb := pubKeyObject.VerifyBytes(s256[:32], pDec)
	if !bb {
		fmt.Println("Bad signature, transaction failed.", sigS)
		return types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("Bad signature. Got %v", sigS) }
	}

	appStore.DeleteCertKey(dcS)
	appStore.SetNonce(addrFrom, nonce+1)

	tags := []cmn.KVPair{
		{Key: []byte("app.type"), Value: []byte(txcmm.TxMsgTypeRemoveCertMsg)},
	}

	return types.ResponseDeliverTx{Code: code.CodeTypeOK, GasUsed: 0, Tags: tags}
}
