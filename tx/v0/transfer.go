package v0

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/Ankr-network/ankr-chain/account"
	ankrcmm "github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/common/code"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	"github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

const (
	MIN_TOKEN_SEND = "5000000000000000000" // 5 tokens
)

type transferMsg struct {}

func (tr *transferMsg) ProcessTx(txMsg interface{}, appStore appstore.AppStore) types.ResponseDeliverTx {
	trxSendSlices, ok := txMsg.([]string)
	if !ok {
		return  types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("invalid tx transfer msg") }
	}

	if len(trxSendSlices) < 6 {
		return types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("Expected trx send. Got %v", trxSendSlices) }
	}

	fromS   := trxSendSlices[0]
	toS     := trxSendSlices[1]
	amountS := trxSendSlices[2]
	nonceS  := trxSendSlices[3]
	pubkeyS := trxSendSlices[4]
	sigS    := trxSendSlices[5]

	if len(fromS) != ankrcmm.KeyAddressLen {
		return  types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("Unexpected from address. Got %v", fromS) }
	}

	if len(toS) != ankrcmm.KeyAddressLen {
		return  types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("Unexpected to address. Got %v", toS)}
	}

	amountSend, ret := new(big.Int).SetString(string(amountS), 10)
	if !ret {
		return  types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("Unexpected amount. Got %v", amountS) }
	} else { // amountSend < 0 or less than MIN_TOKEN_SEND
		minN, _ := new(big.Int).SetString(MIN_TOKEN_SEND, 10)
		if amountSend.Cmp(minN) == -1 || amountSend.Cmp(minN) == 0 {
			return  types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("Unexpected amount, not enough amount. Got %v", amountS) }
		}

		zeroN, _ := new(big.Int).SetString("0", 10)
		if amountSend.Cmp(zeroN) == -1 {
			return types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("Unexpected amount, negative num. Got %v", amountS) }
		}
	}

	nonceInt, err_nonce := strconv.ParseUint(string(nonceS), 10, 64)
	if err_nonce != nil {
		return types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("Unexpected nonce. Got %v", nonceS) }
	}

	nonce, _, _, _, _ := appStore.Nonce(fromS, 0, false)
	if nonceInt != nonce + 1 && nonceInt != nonce {
		return types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("Unexpected nonce. fromS %v, Got %v, Expected %v", fromS, nonceS,nonce ) }
	}

	if len(pubkeyS) == 0 {
		return types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("Unexpected public key. Got %v", pubkeyS) }
	}

	if len(sigS) == 0 {
		return types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("Unexpected signature. Got %v", sigS) }
	}

	// ensure pubkey match fromaddress, you can't send other person's money.
	addr, err := ankrcmm.AddressByPublicKey(string(pubkeyS))
	if err != nil {
		return types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("Parse address error. Got %v", pubkeyS) }
	}

	if string(fromS) != string(addr) {
		return types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("FromAddress no match with pubkey. Got %v", pubkeyS) }
	}

	pDec, _ := base64.StdEncoding.DecodeString(sigS)

	pubKeyObject, err := ankrcmm.DeserilizePubKey(pubkeyS)
	if err != nil {
		return types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("Deserilize pubkey failure. Got %v", pubkeyS) }
	}

	s256 := ankrcmm.ConvertBySha256([]byte(fmt.Sprintf("%s%s%s%s", fromS, toS, amountS, nonceS)))
	bb := pubKeyObject.VerifyBytes(s256[:32], pDec)
	if !bb {
		fmt.Println("Bad signature, transaction failed.", sigS)
		return  types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("Bad signature. Got %v", sigS) }
	}

	fromBalanceInt, _, _, _, err := appStore.Balance(fromS, "ANKR", 0, false)
	if err != nil {
		return  types.ResponseDeliverTx{ Code: code.CodeTypeLoadBalError, Log: fmt.Sprintf("can't get balance, address %v", fromS) }
	}

	zeroN, _ := new(big.Int).SetString("0", 10)
	if fromBalanceInt.Cmp(zeroN) == -1 {
		return types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("Unexpected amount, negative num. Got %v", fromBalanceInt) }
	}

	if fromBalanceInt.Cmp(amountSend) == -1 {
		return types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("Balance not enough. Balance %v, needed amount %v", fromBalanceInt, amountSend) }
	}

	fundBalanceInt, _, _, _, err := appStore.Balance(account.AccountManagerInstance().FoundAccountAddress(), "ANKR", 0, false)
	if err != nil {
		return  types.ResponseDeliverTx{ Code: code.CodeTypeLoadBalError, Log: fmt.Sprintf("can't get balance, address %v", account.AccountManagerInstance().FoundAccountAddress()) }
	}

	if fundBalanceInt.Cmp(zeroN) == -1 {
		return types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("Unexpected amount, negative num. Got %s", fundBalanceInt.String()) }
	}

	toBalanceInt, _, _, _, err := appStore.Balance(toS, "ANKR", 0, false)
	if err != nil {
		toBalanceInt = new(big.Int).SetUint64(0)
	}

	if toBalanceInt.Cmp(zeroN) == -1 {
		return types.ResponseDeliverTx{ Code: code.CodeTypeEncodingError, Log: fmt.Sprintf("Unexpected amount, negative num. Got %v", toBalanceInt.String()) }
	}

	gas, _ := new(big.Int).SetString(MIN_TOKEN_SEND, 10)
	fromBalanceInt.Sub(fromBalanceInt, amountSend)
	// 1. calculate gas based on amountSend
	// 2. actualAmountSend = (amountSend - gas)
	toBalanceInt.Add(toBalanceInt, amountSend.Sub(amountSend, gas))
	fundBalanceInt.Add(fundBalanceInt, gas)

	appStore.SetBalance(fromS, ankrcmm.Amount{ankrcmm.Currency{"ANKR", 18}, fromBalanceInt.Bytes()})
	appStore.SetBalance(toS, ankrcmm.Amount{ankrcmm.Currency{"ANKR", 18}, toBalanceInt.Bytes()}) // use original nonce
	appStore.SetBalance(account.AccountManagerInstance().FoundAccountAddress(), ankrcmm.Amount{ankrcmm.Currency{"ANKR", 18}, fundBalanceInt.Bytes()}) // use original nonce

	if fromS != toS  &&  nonceInt != 1 {
		appStore.SetNonce(fromS, nonce+1)
	}

	tvalue := time.Now().UnixNano()
	addressIndexFrom := fmt.Sprintf("app.%s", fromS)
	addressIndexTo   := fmt.Sprintf("app.%s", toS)
	tags := []cmn.KVPair{
		{Key: []byte("app.fromaddress"), Value: []byte(fromS)},
		{Key: []byte("app.toaddress"), Value: []byte(toS)},
		{Key: []byte(addressIndexFrom), Value: []byte(strconv.FormatInt(1, 10))},
		{Key: []byte(addressIndexTo), Value: []byte(strconv.FormatInt(1, 10))},
		{Key: []byte("app.timestamp"), Value: []byte(strconv.FormatInt(tvalue, 10))},
		{Key: []byte("app.type"), Value: []byte("Send")},
	}

	gasUsed, _ := strconv.ParseInt(MIN_TOKEN_SEND, 0, 64)

	return types.ResponseDeliverTx{Code: code.CodeTypeOK, GasUsed: gasUsed, Tags: tags}
}
