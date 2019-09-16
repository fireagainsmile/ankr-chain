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
	ankrcrypto "github.com/Ankr-network/ankr-chain/crypto"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	"github.com/Ankr-network/ankr-chain/tx"
	"github.com/Ankr-network/ankr-chain/tx/account"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

const (
	MIN_TOKEN_SEND = "5000000000000000000" // 5 tokens
)

func NewTransferTxM() *tx.TxMsg {
	return &tx.TxMsg{ImplTxMsg: new(TransferMsg)}
}

type TransferMsg struct {
}


func (tr *TransferMsg) GasWanted() int64 {
	return 0
}

func (tr *TransferMsg) GasUsed() int64 {
	gasUsed, _ := strconv.ParseInt(MIN_TOKEN_SEND, 0, 64)

	return gasUsed
}

func (tr *TransferMsg) Type() string {
	return ankrtypes.TrxSendPrefix
}

func (tr *TransferMsg) Bytes() []byte {
	return nil
}
func (tr *TransferMsg) SetSecretKey(sk ankrcrypto.SecretKey) {

}

func (tr *TransferMsg) SecretKey() ankrcrypto.SecretKey {
	return nil
}

func (tr *TransferMsg) ProcessTx(txMsg interface{}, appStore appstore.AppStore, isOnlyCheck bool) (uint32, string, []cmn.KVPair){
	trxSendSlices, ok := txMsg.([]string)
	if !ok {
		return  code.CodeTypeEncodingError, fmt.Sprintf("invalid tx transfer msg"), nil
	}

	if len(trxSendSlices) < 6 {
		return code.CodeTypeEncodingError, fmt.Sprintf("Expected trx send. Got %v", trxSendSlices), nil
	}

	fromS   := trxSendSlices[0]
	toS     := trxSendSlices[1]
	amountS := trxSendSlices[2]
	nonceS  := trxSendSlices[3]
	pubkeyS := trxSendSlices[4]
	sigS    := trxSendSlices[5]
	//fmt.Println(fromS, toS, amountS, nonceS, pubkeyS,  sigS)

	if len(fromS) != ankrtypes.KeyAddressLen {
		return  code.CodeTypeEncodingError, fmt.Sprintf("Unexpected from address. Got %v", fromS), nil
	}

	if len(toS) != ankrtypes.KeyAddressLen {
		return  code.CodeTypeEncodingError, fmt.Sprintf("Unexpected to address. Got %v", toS), nil
	}

	amountSend, ret := new(big.Int).SetString(string(amountS), 10)
	if !ret {
		return  code.CodeTypeEncodingError, fmt.Sprintf("Unexpected amount. Got %v", amountS), nil
	} else { // amountSend < 0 or less than MIN_TOKEN_SEND
		minN, _ := new(big.Int).SetString(MIN_TOKEN_SEND, 10)
		if amountSend.Cmp(minN) == -1 || amountSend.Cmp(minN) == 0 {
			return  code.CodeTypeEncodingError, fmt.Sprintf("Unexpected amount, not enough amount. Got %v", amountS), nil
		}

		zeroN, _ := new(big.Int).SetString("0", 10)
		if amountSend.Cmp(zeroN) == -1 {
			return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected amount, negative num. Got %v", amountS), nil
		}
	}

	nonceInt, err_nonce := strconv.ParseInt(string(nonceS), 10, 64)
	if err_nonce != nil {
		return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected nonce4. Got %v", nonceS), nil
	}

	if len(pubkeyS) == 0 {
		return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected public key. Got %v", pubkeyS), nil
	}

	if len(sigS) == 0 {
		return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected signature. Got %v", sigS), nil
	}

	// ensure pubkey match fromaddress, you can't send other person's money.
	addr, err := common.AddressByPublicKey(string(pubkeyS))
	if err != nil {
		return code.CodeTypeEncodingError, fmt.Sprintf("Parse address error. Got %v", pubkeyS), nil
	}

	if string(fromS) != string(addr) {
		return code.CodeTypeEncodingError, fmt.Sprintf("FromAddress no match with pubkey. Got %v", pubkeyS), nil
	}

	pDec, _ := base64.StdEncoding.DecodeString(sigS)

	pubKeyObject, err := common.DeserilizePubKey(pubkeyS)
	if err != nil {
		return code.CodeTypeEncodingError, fmt.Sprintf("Deserilize pubkey failure. Got %v", pubkeyS), nil
	}

	s256 := common.ConvertBySha256([]byte(fmt.Sprintf("%s%s%s%s", fromS, toS, amountS, nonceS)))
	bb := pubKeyObject.VerifyBytes(s256[:32], pDec)
	if !bb {
		fmt.Println("Bad signature, transaction failed.", sigS)
		return  code.CodeTypeEncodingError, fmt.Sprintf("Bad signature. Got %v", sigS), nil
	}

	fromBalanceNonce := appStore.Balance(ankrtypes.PrefixBalanceKey([]byte(fromS)))
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
		return code.CodeTypeEncodingError, fmt.Sprintf("Expected from balance and nonce. Got %v", balanceNonceSlices), nil
	}

	fromBalanceInt, ret := new(big.Int).SetString(string(fromBalance), 10)
	if !ret {
		return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected amount. Got %v", fromBalance), nil
	} else { // fromBalanceInt < 0
		zeroN, _ := new(big.Int).SetString("0", 10)
		if fromBalanceInt.Cmp(zeroN) == -1 {
			return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected amount, negative num. Got %v", fromBalance), nil
		}
	}

	// check stake here. If from balance is less than stake, let it fail.
	var isSucess bool
	cstake, _ := new(big.Int).SetString("0", 10)
	value := appStore.Get([]byte(ankrtypes.AccountStakePrefix))
	if value == nil || string(value) == "" {
		// do nothing for now
	} else {
		stakeNonceSlices := strings.Split(string(value), ":")
		cstake, isSucess = new(big.Int).SetString(string(stakeNonceSlices[0]), 10)
		if !isSucess {
			return code.CodeTypeEncodingError, fmt.Sprintf("stake format error, %v", stakeNonceSlices[0]), nil
		}
	}

	amountSendTemp, _ := new(big.Int).SetString(amountSend.String(), 10)
	amountSendTemp.Add(amountSendTemp, cstake)
	if fromBalanceInt.Cmp(amountSendTemp) == -1 {
		return code.CodeTypeEncodingError, fmt.Sprintf("Balance not enough. Balance %v, needed amount %v", fromBalanceInt, amountSendTemp), nil
	}

	/* check nonce */
	fromNonceInt, err_from := strconv.ParseInt(string(fromNonce), 10, 64)
	if err_from != nil || fromNonceInt < 0 {
		return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected from nonce. Got %v", fromNonce), nil
	}

	if (len(balanceNonceSlices) == 2) && (fromNonceInt + 1) != nonceInt {
		return code.CodeTypeEncodingError, fmt.Sprintf("nonce should be one more than last nonce. Got %v", nonceS), nil
	}

	fundRealBalance, _ := new(big.Int).SetString("0", 10)
	fundBalanceNonce := appStore.Balance(ankrtypes.PrefixBalanceKey([]byte(account.AccountManagerInstance().GasAccountAddress())))
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
			return code.CodeTypeEncodingError, fmt.Sprintf("Expected to balance and nonce of fund. Got %v", balanceNonceSlices), nil
		}
	}

	if fundBalanceNonce != nil {
		fundBalanceInt, err := new(big.Int).SetString(string(fundBalance), 10)
		if !err {
			return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected amount. Got %v", fundBalance), nil
		} else { // toBalanceInt < 0
			zeroN, _ := new(big.Int).SetString("0", 10)
			if fundBalanceInt.Cmp(zeroN) == -1 {
				return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected amount, negative num. Got %v", fundBalance), nil
			}
		}

		fundRealBalance = fundBalanceInt
	}

	toRealBalance, _ := new(big.Int).SetString("0", 10)
	toBalanceNonce := appStore.Balance(ankrtypes.PrefixBalanceKey([]byte(toS)))
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
			return code.CodeTypeEncodingError, fmt.Sprintf("Expected to balance and nonce. Got %v", balanceNonceSlices), nil
		}
	}

	if toBalanceNonce != nil {
		toBalanceInt, err := new(big.Int).SetString(string(toBalance), 10)
		if !err {
			return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected amount. Got %v", toBalance), nil
		} else { // toBalanceInt < 0
			zeroN, _ := new(big.Int).SetString("0", 10)
			if toBalanceInt.Cmp(zeroN) == -1 {
				return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected amount, negative num. Got %v", toBalance), nil
			}
		}

		toRealBalance = toBalanceInt
	}

	if isOnlyCheck {
		return code.CodeTypeOK, "", nil
	}

	gas, _ := new(big.Int).SetString(MIN_TOKEN_SEND, 10)
	fromBalanceInt.Sub(fromBalanceInt, amountSend)
	// 1. calculate gas based on amountSend
	// 2. actualAmountSend = (amountSend - gas)
	toRealBalance.Add(toRealBalance, amountSend.Sub(amountSend, gas))
	fundRealBalance.Add(fundRealBalance, gas)

	appStore.SetBalance(ankrtypes.PrefixBalanceKey([]byte(fromS)), []byte(fromBalanceInt.String()+":"+nonceS))
	appStore.SetBalance(ankrtypes.PrefixBalanceKey([]byte(toS)), []byte(toRealBalance.String()+":"+toNonce)) // use original nonce
	appStore.SetBalance(ankrtypes.PrefixBalanceKey([]byte(account.AccountManagerInstance().GasAccountAddress())), []byte(fundRealBalance.String()+":"+fundNonce)) // use original nonce

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

	return code.CodeTypeOK, "", tags
}
