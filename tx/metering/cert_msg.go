package metering


import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/Ankr-network/ankr-chain/account"
	ankrcmm "github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/common/code"
	ankrcrypto "github.com/Ankr-network/ankr-chain/crypto"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	"github.com/Ankr-network/ankr-chain/tx"
	txcmm "github.com/Ankr-network/ankr-chain/tx/common"
	"github.com/go-interpreter/wagon/exec/gas"
	cmn "github.com/tendermint/tendermint/libs/common"
)

type SetCertMsg struct {
	FromAddr  string  `json:"fromaddr"`
	DCName    string  `json:"dcname"`
	NSName    string  `json:"nsname"`
	PemBase64 string  `json:"pembase64"`
}

func (sc *SetCertMsg) SignerAddr() []string {
	return []string {sc.FromAddr}
}

func (sc *SetCertMsg) Type() string {
	return txcmm.TxMsgTypeSetCertMsg
}

func (sc *SetCertMsg) Bytes(txSerializer tx.TxSerializer) []byte {
	bytesRtn, _ := txSerializer.MarshalJSON(sc)

	return bytesRtn
}

func (sc *SetCertMsg) SetSecretKey(sk ankrcrypto.SecretKey) {

}

func (sc *SetCertMsg) SecretKey() ankrcrypto.SecretKey {
	return &ankrcrypto.SecretKeyEd25519{}
}

func (sc *SetCertMsg) PermitKey(store appstore.AppStore, pubKey []byte) bool {
	adminPubkey := store.Get([]byte(ankrcmm.ADMIN_OP_METERING_PUBKEY_NAME))
	if len(adminPubkey) == 0 {
		adminPubkey = []byte(account.AccountManagerInstance().AdminOpAccount(ankrcmm.AccountAdminMetering))
	}

	adminPubKeyStr, err := base64.StdEncoding.DecodeString(string(adminPubkey))
	if err != nil {
		return false
	}

	return  bytes.Equal(pubKey, []byte(adminPubKeyStr))
}

func (sc *SetCertMsg) ProcessTx(context tx.ContextTx, metric gas.GasMetric, isOnlyCheck bool) (uint32, string, []cmn.KVPair) {
	if len(sc.FromAddr) != ankrcmm.KeyAddressLen {
		return  code.CodeTypeInvalidAddress, fmt.Sprintf("SetCertMsg ProcessTx, unexpected from address. Got %s, addr len=%d", sc.FromAddr, len(sc.FromAddr)), nil
	}

	if isOnlyCheck {
		return code.CodeTypeOK, "", nil
	}

	context.AppStore().SetCertKey(sc.DCName, sc.NSName, sc.PemBase64)

	context.AppStore().IncNonce(sc.FromAddr)

	tags := []cmn.KVPair{
		{Key: []byte("app.type"), Value: []byte(txcmm.TxMsgTypeSetCertMsg)},
	}

	return code.CodeTypeOK, "", tags
}

type RemoveCertMsg struct {
	FromAddr  string  `json:"fromaddr"`
	DCName    string  `json:"dcname"`
	NSName    string  `json:"nsname"`
}

func (rc *RemoveCertMsg) SignerAddr() []string {
	return []string {rc.FromAddr}
}

func (rc *RemoveCertMsg) GasWanted() int64 {
	return 0
}

func (rc *RemoveCertMsg) GasUsed() int64 {
	return 0
}

func (rc *RemoveCertMsg) Type() string {
	return txcmm.TxMsgTypeRemoveCertMsg
}

func (rc *RemoveCertMsg) Bytes(txSerializer tx.TxSerializer) []byte {
	bytesRtn, _ := txSerializer.MarshalJSON(rc)

	return bytesRtn
}

func (rc *RemoveCertMsg) SetSecretKey(sk ankrcrypto.SecretKey) {

}

func (rc *RemoveCertMsg) SecretKey() ankrcrypto.SecretKey {
	return &ankrcrypto.SecretKeyEd25519{}
}

func (sc *RemoveCertMsg) PermitKey(store appstore.AppStore, pubKey []byte) bool {
	adminPubkey := store.Get([]byte(ankrcmm.ADMIN_OP_METERING_PUBKEY_NAME))
	if len(adminPubkey) == 0 {
		adminPubkey = []byte(account.AccountManagerInstance().AdminOpAccount(ankrcmm.AccountAdminMetering))
	}

	adminPubKeyStr, err := base64.StdEncoding.DecodeString(string(adminPubkey))
	if err != nil {
		return false
	}

	return  bytes.Equal(pubKey, []byte(adminPubKeyStr))
}

func (rc *RemoveCertMsg) ProcessTx(context tx.ContextTx, metric gas.GasMetric, isOnlyCheck bool) (uint32, string, []cmn.KVPair) {
	if len(rc.FromAddr) != ankrcmm.KeyAddressLen {
		return  code.CodeTypeInvalidAddress, fmt.Sprintf("RemoveCertMsg ProcessTx, unexpected from address. Got %s, addr len=%d", rc.FromAddr, len(rc.FromAddr)), nil
	}

	if isOnlyCheck {
		return code.CodeTypeOK, "", nil
	}

	context.AppStore().DeleteCertKey(rc.DCName, rc.NSName)

	context.AppStore().IncNonce(rc.FromAddr)

	tags := []cmn.KVPair{
		{Key: []byte("app.type"), Value: []byte(txcmm.TxMsgTypeRemoveCertMsg)},
	}

	return code.CodeTypeOK, "", tags
}

