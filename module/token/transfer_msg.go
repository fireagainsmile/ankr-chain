package token

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/common/code"
	"github.com/Ankr-network/ankr-chain/module/account"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	"github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

const (
	MIN_TOKEN_SEND = "5000000000000000000" // 5 tokens
)



type TransferMsg struct {
}

func (tr *TransferMsg) CheckTx(tx []byte, appStore appstore.AppStore) types.ResponseCheckTx {
	tx = tx[len(ankrtypes.TrxSendPrefix):]
	trxSendSlices := strings.Split(string(tx), ":")
	if len(trxSendSlices) < 6 {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Transaction send incorecct format, got %s", string(tx))}
	}

	fromS := trxSendSlices[0]
	toS := trxSendSlices[1]
	amountS := trxSendSlices[2]
	nonceS := trxSendSlices[3]

	if len(fromS) != ankrtypes.KeyAddressLen {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected from address. Got %v", fromS)}
	}

	if len(toS) != ankrtypes.KeyAddressLen {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected to address. Got %v", toS)}
	}

	amountSend, err := new(big.Int).SetString(string(amountS), 10)
	if !err {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected amount. Got %v", amountS)}
	} else { // amountSend < 0 or less than min
		zeroN, _ := new(big.Int).SetString("0", 10)
		if amountSend.Cmp(zeroN) == -1 {
			return types.ResponseCheckTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Unexpected amount, negative num. Got %v", amountS)}
		}

		minN, _ := new(big.Int).SetString(MIN_TOKEN_SEND, 10)
		if amountSend.Cmp(minN) == -1 || amountSend.Cmp(minN) == 0 {
			return types.ResponseCheckTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Unexpected amount, not enough amount. Got %v", amountS)}
		}
	}

	nonceInt, err_nonce := strconv.ParseInt(string(nonceS), 10, 64)
	if err_nonce != nil {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected nonce. Got %v", nonceS)}
	}


	fromBalanceNonce := appStore.Get(ankrtypes.PrefixBalanceKey([]byte(fromS)))
	balanceNonceSlices := strings.Split(string(fromBalanceNonce), ":")
	var fromBalance string
	var fromNonce string
	if len(balanceNonceSlices) == 1 {
		fromBalance = balanceNonceSlices[0]
		fromNonce = "1"
	} else if len(balanceNonceSlices) == 2 {
		fromBalance = balanceNonceSlices[0]
		fromNonce = balanceNonceSlices[1]
	} else {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Expected from balance and nonce. Got %v", balanceNonceSlices)}
	}

	fromBalanceInt, err := new(big.Int).SetString(string(fromBalance), 10)
	if !err {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected amount. Got %v", fromBalance)}
	} else { // fromBalanceInt < 0
		zeroN, _ := new(big.Int).SetString("0", 10)
		if fromBalanceInt.Cmp(zeroN) == -1 {
			return types.ResponseCheckTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Unexpected amount, negative num. Got %v", fromBalance)}
		}
	}

	if fromBalanceInt.Cmp(amountSend) == -1 { // bignumber comparison
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Not enough balance to send. Balance %v, send %v", fromBalanceInt, amountSend)}
	}

	// check stake here. If from balance is less than stake, let it fail.
	cstake, _ := new(big.Int).SetString("0", 10)
	value := appStore.Get([]byte(ankrtypes.AccountStakePrefix))
	if value == nil || string(value) == "" {
		// do nothing for now
	} else {
		stakeNonceSlices := strings.Split(string(value), ":")
		cstake, err = new(big.Int).SetString(string(stakeNonceSlices[0]), 10)
		if !err {
			return types.ResponseCheckTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("stake format error, %v", stakeNonceSlices[0])}
		}
	}

	if fromBalanceInt.Cmp(cstake) == -1 {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Balance <= stake. Balance %v, stake %v", fromBalanceInt, cstake)}
	}

	// check stake again.
	fromBalanceInt.Sub(fromBalanceInt, amountSend)
	if fromBalanceInt.Cmp(cstake) == -1 {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Sub Balance <= stake. Balance %v, stake %v", fromBalanceInt, cstake)}
	}


	/* check nonce */
	fromNonceInt, err_from := strconv.ParseInt(string(fromNonce), 10, 64)
	if err_from != nil || fromNonceInt < 0 {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected from nonce. Got %v", fromNonce)}
	}

	if (len(balanceNonceSlices) == 2) && (fromNonceInt + 1) != nonceInt {
		return types.ResponseCheckTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("nonce should be one more than last nonce. Got %v", nonceS)}
	}

	return types.ResponseCheckTx{Code: code.CodeTypeOK, GasWanted: 1}
}

// format is "trx_send=from:to:amount:nonce:pubkey:sig"
// nonce should be stored in from account.
// will add signature verification when wallet code is ready
func (tr *TransferMsg) DeliverTx(tx []byte, appStore appstore.AppStore) types.ResponseDeliverTx {
	tx = tx[len(ankrtypes.TrxSendPrefix):]
	trxSendSlices := strings.Split(string(tx), ":")
	if len(trxSendSlices) < 6 {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Expected trx send. Got %v", trxSendSlices)}
	}

	fromS := trxSendSlices[0]
	toS := trxSendSlices[1]
	amountS := trxSendSlices[2]
	nonceS := trxSendSlices[3]
	pubkeyS := trxSendSlices[4]
	sigS := trxSendSlices[5]
	//fmt.Println(fromS, toS, amountS, nonceS, pubkeyS,  sigS)

	if len(fromS) != ankrtypes.KeyAddressLen {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected from address. Got %v", fromS)}
	}

	if len(toS) != ankrtypes.KeyAddressLen {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected to address. Got %v", toS)}
	}

	amountSend, ret := new(big.Int).SetString(string(amountS), 10)
	if !ret {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected amount. Got %v", amountS)}
	} else { // amountSend < 0 or less than MIN_TOKEN_SEND
		minN, _ := new(big.Int).SetString(MIN_TOKEN_SEND, 10)
		if amountSend.Cmp(minN) == -1 || amountSend.Cmp(minN) == 0 {
			return types.ResponseDeliverTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Unexpected amount, not enough amount. Got %v", amountS)}
		}

		zeroN, _ := new(big.Int).SetString("0", 10)
		if amountSend.Cmp(zeroN) == -1 {
			return types.ResponseDeliverTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Unexpected amount, negative num. Got %v", amountS)}
		}
	}

	nonceInt, err_nonce := strconv.ParseInt(string(nonceS), 10, 64)
	if err_nonce != nil {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected nonce4. Got %v", nonceS)}
	}

	if len(pubkeyS) == 0 {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected public key. Got %v", pubkeyS)}
	}

	if len(sigS) == 0 {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected signature. Got %v", sigS)}
	}

	// ensure pubkey match fromaddress, you can't send other person's money.
	addr, err := common.AddressByPublicKey(string(pubkeyS))
	if err != nil {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Parse address error. Got %v", pubkeyS)}
	}

	if string(fromS) != string(addr) {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("FromAddress no match with pubkey. Got %v", pubkeyS)}
	}

	pDec, _ := base64.StdEncoding.DecodeString(sigS)

	pubKeyObject, err := common.DeserilizePubKey(pubkeyS)
	if err != nil {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Deserilize pubkey failure. Got %v", pubkeyS)}
	}

	s256 := common.ConvertBySha256([]byte(fmt.Sprintf("%s%s%s%s", fromS, toS, amountS, nonceS)))
	bb := pubKeyObject.VerifyBytes(s256[:32], pDec)
	if !bb {
		fmt.Println("Bad signature, transaction failed.", sigS)
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Bad signature. Got %v", sigS)}
	}

	fromBalanceNonce := appStore.Get(ankrtypes.PrefixBalanceKey([]byte(fromS)))
	balanceNonceSlices := strings.Split(string(fromBalanceNonce), ":")
	var fromBalance string
	var fromNonce string
	if len(balanceNonceSlices) == 1 {
		fromBalance = balanceNonceSlices[0]
		fromNonce = "1"
	} else if len(balanceNonceSlices) == 2 {
		fromBalance = balanceNonceSlices[0]
		fromNonce = balanceNonceSlices[1]
	} else {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Expected from balance and nonce. Got %v", balanceNonceSlices)}
	}

	fromBalanceInt, ret := new(big.Int).SetString(string(fromBalance), 10)
	if !ret {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected amount. Got %v", fromBalance)}
	} else { // fromBalanceInt < 0
		zeroN, _ := new(big.Int).SetString("0", 10)
		if fromBalanceInt.Cmp(zeroN) == -1 {
			return types.ResponseDeliverTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Unexpected amount, negative num. Got %v", fromBalance)}
		}
	}

	if fromBalanceInt.Cmp(amountSend) == -1 { // bignumber comparison
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Not enough balance to send. Balance %v, send %v", fromBalanceInt, amountSend)}
	}

	/* check nonce */
	fromNonceInt, err_from := strconv.ParseInt(string(fromNonce), 10, 64)
	if err_from != nil || fromNonceInt < 0 {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("Unexpected from nonce. Got %v", fromNonce)}
	}

	if (len(balanceNonceSlices) == 2) && (fromNonceInt + 1) != nonceInt {
		return types.ResponseDeliverTx{
			Code: code.CodeTypeEncodingError,
			Log:  fmt.Sprintf("nonce should be one more than last nonce. Got %v", nonceS)}
	}

	fundRealBalance, _ := new(big.Int).SetString("0", 10)
	fundBalanceNonce := appStore.Get(ankrtypes.PrefixBalanceKey([]byte(account.AccountManagerInstance().GasAccountAddress())))
	var fundBalance string
	var fundNonce string = "1"
	if  fundBalanceNonce != nil {
		balanceNonceSlices = strings.Split(string(fundBalanceNonce), ":")
		if (len(balanceNonceSlices) == 1) {
			fundBalance = balanceNonceSlices[0]
			fundNonce = "1"
		} else if len(balanceNonceSlices) == 2 {
			fundBalance = balanceNonceSlices[0]
			fundNonce = balanceNonceSlices[1]
			if fundNonce == "" {
				fundNonce = "1"
			}
		} else {
			return types.ResponseDeliverTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Expected to balance and nonce of fund. Got %v", balanceNonceSlices)}
		}
	}

	if fundBalanceNonce != nil {
		fundBalanceInt, err := new(big.Int).SetString(string(fundBalance), 10)
		if !err {
			return types.ResponseDeliverTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Unexpected amount. Got %v", fundBalance)}
		} else { // toBalanceInt < 0
			zeroN, _ := new(big.Int).SetString("0", 10)
			if fundBalanceInt.Cmp(zeroN) == -1 {
				return types.ResponseDeliverTx{
					Code: code.CodeTypeEncodingError,
					Log:  fmt.Sprintf("Unexpected amount, negative num. Got %v", fundBalance)}
			}
		}

		fundRealBalance = fundBalanceInt
	}

	toRealBalance, _ := new(big.Int).SetString("0", 10)
	toBalanceNonce := appStore.Get(ankrtypes.PrefixBalanceKey([]byte(toS)))
	var toBalance string
	var toNonce string = "1"
	if  toBalanceNonce != nil {
		balanceNonceSlices = strings.Split(string(toBalanceNonce), ":")
		if (len(balanceNonceSlices) == 1) {
			toBalance = balanceNonceSlices[0]
			toNonce = "1"
		} else if len(balanceNonceSlices) == 2 {
			toBalance = balanceNonceSlices[0]
			toNonce = balanceNonceSlices[1]
			if toNonce == "" {
				toNonce = "1"
			}
		} else {
			return types.ResponseDeliverTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Expected to balance and nonce. Got %v", balanceNonceSlices)}
		}
	}

	if toBalanceNonce != nil {
		toBalanceInt, err := new(big.Int).SetString(string(toBalance), 10)
		if !err {
			return types.ResponseDeliverTx{
				Code: code.CodeTypeEncodingError,
				Log:  fmt.Sprintf("Unexpected amount. Got %v", toBalance)}
		} else { // toBalanceInt < 0
			zeroN, _ := new(big.Int).SetString("0", 10)
			if toBalanceInt.Cmp(zeroN) == -1 {
				return types.ResponseDeliverTx{
					Code: code.CodeTypeEncodingError,
					Log:  fmt.Sprintf("Unexpected amount, negative num. Got %v", toBalance)}
			}
		}

		toRealBalance = toBalanceInt
	}

	//fmt.Println(toRealBalance)

	gas, _ := new(big.Int).SetString(MIN_TOKEN_SEND, 10)
	fromBalanceInt.Sub(fromBalanceInt, amountSend)
	// 1. calculate gas based on amountSend
	// 2. actualAmountSend = (amountSend - gas)
	toRealBalance.Add(toRealBalance, amountSend.Sub(amountSend, gas))
	fundRealBalance.Add(fundRealBalance, gas)

	appStore.Set(ankrtypes.PrefixBalanceKey([]byte(fromS)), []byte(fromBalanceInt.String()+":"+nonceS))
	appStore.Set(ankrtypes.PrefixBalanceKey([]byte(toS)), []byte(toRealBalance.String()+":"+toNonce)) // use original nonce
	appStore.Set(ankrtypes.PrefixBalanceKey([]byte(account.AccountManagerInstance().GasAccountAddress())), []byte(fundRealBalance.String()+":"+fundNonce)) // use original nonce
	appStore.IncSize()

	tvalue := time.Now().UnixNano()
	tags := []cmn.KVPair{
		{Key: []byte("app.fromaddress"), Value: []byte(fromS)},
		{Key: []byte("app.toaddress"), Value: []byte(toS)},
		{Key: []byte("app.timestamp"), Value: []byte(strconv.FormatInt(tvalue, 10))},
		{Key: []byte("app.type"), Value: []byte("Send")},
	}

	gasUsed, _ := strconv.ParseInt(MIN_TOKEN_SEND, 0, 64)
	return types.ResponseDeliverTx{Code: code.CodeTypeOK, GasUsed: gasUsed, Tags: tags}
}
