package tx

import (

	"github.com/Ankr-network/ankr-chain/account"

	"github.com/Ankr-network/ankr-chain/common/code"
	ankrcrypto "github.com/Ankr-network/ankr-chain/crypto"
	"github.com/Ankr-network/ankr-chain/store/appstore"


	cmn "github.com/tendermint/tendermint/libs/common"
)

type TxMsgTesting struct {
	ToAddr  string            `json:"toaddr"`
	Asserts []account.Amount `json:"asserts"`
}

func (tr *TxMsgTesting) FromAddress() string {
	return ""
}

func (tr *TxMsgTesting) GasWanted() int64 {
	return 0
}

func (tr *TxMsgTesting) GasUsed() int64 {

	return 0
}

func (tr *TxMsgTesting) Type() string {
	return "Testing"
}

func (tr *TxMsgTesting) Bytes() []byte {
	return nil
}
func (tr *TxMsgTesting) SetSecretKey(sk ankrcrypto.SecretKey) {

}

func (tr *TxMsgTesting) SecretKey() ankrcrypto.SecretKey {
	return nil
}

func (tr *TxMsgTesting) ProcessTx(appStore appstore.AppStore, isOnlyCheck bool) (uint32, string, []cmn.KVPair){

	tags := []cmn.KVPair{
		{Key: []byte("app.fromaddress"), Value: []byte("Address1")},
		{Key: []byte("app.toaddress"), Value: []byte("Address2")},
		{Key: []byte("app.type"), Value: []byte("Send")},
	}

	return code.CodeTypeOK, "", tags
}
