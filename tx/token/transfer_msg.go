package token

import (
	"fmt"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	"math/big"
	"strconv"
	"time"

	"github.com/Ankr-network/ankr-chain/account"
	"github.com/Ankr-network/ankr-chain/common/code"
	ankrcontext "github.com/Ankr-network/ankr-chain/context"
	"github.com/Ankr-network/ankr-chain/contract"
	ankrcrypto "github.com/Ankr-network/ankr-chain/crypto"
	"github.com/Ankr-network/ankr-chain/tx"
	txcmm "github.com/Ankr-network/ankr-chain/tx/common"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

const (
	MIN_TOKEN_SEND = "5000000000000000000" // 5 tokens
)

type TransferMsg struct {
	FromAddr string           `json:"fromaddr"`
	ToAddr   string           `json:"toaddr"`
	Amounts  []account.Amount `json:"amounts"`
	gasUsed  *big.Int
}

type signTransferMsg struct {
	FromAddr string           `json:"fromaddr"`
	ToAddr   string           `json:"toaddr"`
	Amounts  []account.Amount `json:"amounts"`
}

func (stf signTransferMsg) bytes(txSerializer tx.TxSerializer) ([]byte, error) {
	return txSerializer.MarshalJSON(&stf)
}

func (tf *TransferMsg) SignerAddr() []string {
	return []string {tf.FromAddr}
}

func (tf *TransferMsg) GasWanted() int64 {
	return 0
}

func (tf *TransferMsg) GasUsed() int64 {
	gasUsed, _ := strconv.ParseInt(MIN_TOKEN_SEND, 0, 64)

	return gasUsed
}

func (tf *TransferMsg) Type() string {
	return txcmm.TxMsgTypeTransfer
}

func (tf *TransferMsg) signMsg() *signTransferMsg {
	return &signTransferMsg{FromAddr: tf.FromAddr, ToAddr: tf.ToAddr/*, Amounts: tf.Amounts*/}
}

func (tf *TransferMsg) Bytes(txSerializer tx.TxSerializer) []byte {
	bytes, _ :=  tf.signMsg().bytes(txSerializer)
	return bytes
}
func (tf *TransferMsg) SetSecretKey(sk ankrcrypto.SecretKey) {

}

func (tf *TransferMsg) SecretKey() ankrcrypto.SecretKey {
	return &ankrcrypto.SecretKeyEd25519{}
}

func (tf *TransferMsg) PermitKey(store appstore.AppStore, pubKey []byte) bool {
	return true
}

func (tf *TransferMsg) SpendGas(gas *big.Int) bool {
	tf.gasUsed.Add(tf.gasUsed, gas)
	return true
}

func (tf *TransferMsg) SenderAddr() string {
	return tf.FromAddr
}

func (tf *TransferMsg) ProcessTx(context tx.ContextTx, isOnlyCheck bool) (uint32, string, []cmn.KVPair){
	if len(tf.FromAddr) != ankrtypes.KeyAddressLen {
		return  code.CodeTypeInvalidAddress, fmt.Sprintf("TransferMsg ProcessTx, unexpected from address. Got %s, addr len=%d", tf.FromAddr, len(tf.FromAddr)), nil
	}
	if len(tf.ToAddr) != ankrtypes.KeyAddressLen {
		return code.CodeTypeInvalidAddress, fmt.Sprintf("TransferMsg ProcessTx, unexpected to address. Got %s, addr len=%d", tf.ToAddr, len(tf.ToAddr)), nil
	}

	if isOnlyCheck {
		return code.CodeTypeOK, "", nil
	}

	trAmount := tf.Amounts[0]
	tokenContract, err := context.AppStore().LoadContract([]byte(ankrtypes.ContractTokenStorePrefix + trAmount.Cur.Symbol))
	if err != nil {
		return code.CodeTypeLoadContractErr, fmt.Sprintf("load contract err: name = %s", ankrtypes.ContractTokenStorePrefix + trAmount.Cur.Symbol), nil
	}

	params :=  []*ankrtypes.Param{{0, "string", tf.FromAddr},
		{1, "string", tf.ToAddr},
		{2, "string", new(big.Int).SetBytes(tf.Amounts[0].Value).String()},
	}

	contractType    := ankrtypes.ContractType(tokenContract[0])
	contractContext := ankrcontext.NewContextContract(tf, tf, context.AppStore())
    rtn, err := contract.Call(contractContext, contractType, tokenContract[1:], "ANKR", "TransferFrom", params, "bool")
    if err != nil {
    	return code.CodeTypeCallContractErr, fmt.Sprintf("call contract err: contract=%s, method=TransferFrom, err=%v", tf.Amounts[0].Cur.Symbol, err), nil
	}
    isCallSucess := rtn.(bool)
    if !isCallSucess {
		return code.CodeTypeCallContractErr, fmt.Sprintf("call contract err: contract=%s, method=TransferFrom", tf.Amounts[0].Cur.Symbol), nil
	}

	context.AppStore().IncNonce(tf.FromAddr)

	tvalue := time.Now().UnixNano()
	tags := []cmn.KVPair{
		{Key: []byte("app.fromaddress"), Value: []byte(tf.FromAddr)},
		{Key: []byte("app.toaddress"), Value: []byte(tf.ToAddr)},
		{Key: []byte("app.timestamp"), Value: []byte(strconv.FormatInt(tvalue, 10))},
		{Key: []byte("app.type"), Value: []byte(txcmm.TxMsgTypeTransfer)},
	}

	return code.CodeTypeOK, "", tags
}
