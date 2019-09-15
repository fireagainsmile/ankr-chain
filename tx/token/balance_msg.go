package token

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/Ankr-network/ankr-chain/account"
	"github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/common/code"
	ankrcrypto "github.com/Ankr-network/ankr-chain/crypto"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	"github.com/Ankr-network/ankr-chain/tx"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

func NewBalanceTxMsg() *tx.TxMsg {
	return &tx.TxMsg{ImplTxMsg: new(BalanceMsg)}
}

type BalanceMsg struct {
}

func (b *BalanceMsg) GasWanted() int64 {
	return 0
}

func (b *BalanceMsg) GasUsed() int64 {
	return 0
}

func (b *BalanceMsg) Type() string {
	return ankrtypes.SetBalancePrefix
}

func (b *BalanceMsg) Bytes() []byte {
	return nil
}
func (b *BalanceMsg) SetSecretKey(sk ankrcrypto.SecretKey) {

}

func (b *BalanceMsg) SecretKey() ankrcrypto.SecretKey {
	return nil
}

func (b *BalanceMsg) ProcessTx(appStore appstore.AppStore, isOnlyCheck bool) (uint32, string,  []cmn.KVPair) {
	addressS  := trxSetBalanceSlices[0]
	amountS   := trxSetBalanceSlices[1]
	nonceS    := trxSetBalanceSlices[2]
	adminPubS := trxSetBalanceSlices[3]
	sigS      := trxSetBalanceSlices[4]

	var admin_pubkey_str string = ""
	admin_pubkey := appStore.Get([]byte(ankrtypes.ADMIN_OP_FUND_PUBKEY_NAME))
	if len(admin_pubkey) == 0 {
		fmt.Println("use default ADMIN_OP_FUND_PUBKEY_NAME")
		admin_pubkey_str = adminPubkeyOfBalance()
	} else {
		admin_pubkey_str = string(admin_pubkey)
	}

	if adminPubS != admin_pubkey_str {
		return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected pubkey. Got %v", adminPubS), nil
	}

	if len(addressS) != ankrtypes.KeyAddressLen {
		return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected address. Got %v", addressS), nil
	}

	amountSet, err := new(big.Int).SetString(string(amountS), 10)
	if !err {
		return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected amount. Got %v", amountS), nil
	} else { // amountSet < 0
		zeroN, _ := new(big.Int).SetString("0", 10)
		if amountSet.Cmp(zeroN) == -1 {
			return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected amount, negative num. Got %v", amountS), nil
		}
	}

	nonceInt, err_n := strconv.ParseInt(string(nonceS), 10, 64)
	if err_n != nil {
		return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected nonce. Got %v", nonceS), nil
	}

	pDec, _ := base64.StdEncoding.DecodeString(sigS)
	pubKeyObject, err_d := common.DeserilizePubKey(adminPubS)
	if err_d != nil {
		return code.CodeTypeEncodingError, fmt.Sprintf("Deserilize pubkey failure. Got %v", adminPubS), nil
	}

	s256 := common.ConvertBySha256([]byte(fmt.Sprintf("%s%s%s", addressS, amountS, nonceS)))
	bb := pubKeyObject.VerifyBytes(s256[:32], pDec)
	if !bb {
		fmt.Println("Bad signature, transaction failed.", sigS)
		return  code.CodeTypeEncodingError, fmt.Sprintf("Bad signature. Got %v", sigS), nil
	}

	inBalanceAndNonce := appStore.Balance(ankrtypes.PrefixBalanceKey([]byte(addressS)))
	balanceNonceSlices := strings.Split(string(inBalanceAndNonce), ":")
	var inBalance string
	var inNonce string
	if len(balanceNonceSlices) == 1 {
		inBalance = balanceNonceSlices[0]
		inNonce = "0"
	} else if len(balanceNonceSlices) == 2 {
		inBalance = balanceNonceSlices[0]
		inNonce = balanceNonceSlices[1]
	} else {
		inBalance = "0"
		inNonce = "0"
	}
	_ = inBalance

	inNonceInt, err_p:= strconv.ParseInt(string(inNonce), 10, 64)
	if err_p != nil || inNonceInt < 0 {
		return  code.CodeTypeEncodingError, fmt.Sprintf("Unexpected from nonce. Got %v", inNonce), nil
	}

	if (len(balanceNonceSlices) == 2) && (inNonceInt + 1) != nonceInt {
		return  code.CodeTypeEncodingError, fmt.Sprintf("nonce should be one more than last nonce. Got %v", nonceS), nil
	}

	if !isOnlyCheck {
		return code.CodeTypeOK, "", nil
	}

	appStore.SetBalance(ankrtypes.PrefixBalanceKey([]byte(addressS)), []byte(amountS+":"+nonceS))

	tags := []cmn.KVPair{
		{Key: []byte("app.type"), Value: []byte("SetBalance")},
	}

	return code.CodeTypeOK, "", tags
}
