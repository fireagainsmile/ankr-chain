package key

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/Ankr-network/ankr-chain/account"
	"github.com/Ankr-network/ankr-chain/common/code"
	ankrcrypto "github.com/Ankr-network/ankr-chain/crypto"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	"github.com/Ankr-network/ankr-chain/tx"
	txcmm "github.com/Ankr-network/ankr-chain/tx/common"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

func NewKeyMsgTxMsg() *tx.TxMsg {
	return  &tx.TxMsg{ImplTxMsg: new(KeyMsg)}
}

type KeyMsg struct {
	FromAddr string  `json:"fromaddr"`
	KeyName  string  `json:"keyname"`
	KeyValue string  `json:"keyname"`
}

func (k *KeyMsg) SignerAddr() []string {
	return []string {k.FromAddr}
}

func (k *KeyMsg) GasWanted() int64 {
	return 1
}

func (k *KeyMsg) GasUsed() int64 {
	return 0
}

func (s *KeyMsg) Type() string {
	return txcmm.TxMsgTypeKeyMsg
}

func (k *KeyMsg) Bytes(txSerializer tx.TxSerializer) []byte {
	bytes, _ := txSerializer.MarshalJSON(k)

	return bytes
}

func (k *KeyMsg) SetSecretKey(sk ankrcrypto.SecretKey) {

}

func (k *KeyMsg) SecretKey() ankrcrypto.SecretKey {
	return &ankrcrypto.SecretKeyEd25519{}
}

func (k *KeyMsg) PermitKey(store appstore.AppStore, pubKey []byte) bool {
	adminPubkey := account.AccountManagerInstance().AdminOpAccount(account.AccountAdminOP)
	adminPubKeyStr, err := base64.StdEncoding.DecodeString(adminPubkey)
	if err != nil {
		return false
	}

	return  bytes.Equal(pubKey, []byte(adminPubKeyStr))
}

func (k *KeyMsg) ProcessTx(context tx.ContextTx, isOnlyCheck bool) (uint32, string, []cmn.KVPair) {
	if k.KeyName != ankrtypes.ADMIN_OP_VAL_PUBKEY_NAME && k.KeyName != ankrtypes.ADMIN_OP_FUND_PUBKEY_NAME &&
		k.KeyName != ankrtypes.ADMIN_OP_METERING_PUBKEY_NAME {
		return code.CodeTypeEncodingError, fmt.Sprintf("Unexpected keyname. Got %v", k.KeyName), nil
	}

	if isOnlyCheck {
		return  code.CodeTypeOK, "", nil
	}

	context.AppStore().Set([]byte(k.KeyName), []byte(k.KeyValue))

	context.AppStore().IncNonce(k.FromAddr)

	tags := []cmn.KVPair{
		{Key: []byte("app.type"), Value: []byte(txcmm.TxMsgTypeKeyMsg)},
	}

	return code.CodeTypeOK, "", tags
}


