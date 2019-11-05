package metering

import (
	"strconv"
	"time"

	"github.com/Ankr-network/ankr-chain/common/code"
	ankrcrypto "github.com/Ankr-network/ankr-chain/crypto"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	"github.com/Ankr-network/ankr-chain/tx"
	txcmm "github.com/Ankr-network/ankr-chain/tx/common"
	"github.com/Ankr-network/wagon/exec/gas"
	cmn "github.com/tendermint/tendermint/libs/common"
)

func NewMeteringTxMsg() *tx.TxMsg {
	return &tx.TxMsg{ImplTxMsg: new(MeteringMsg)}
}

type MeteringMsg struct {
	FromAddr string  `json:"fromaddr"`
	DCName   string  `json:"dcname"`
	NSName   string  `json:"nsname"`
	Value    string  `json:"value"`
}

func (m *MeteringMsg) SignerAddr() []string {
	return []string {m.FromAddr}
}

func (m *MeteringMsg) Type() string {
	return txcmm.TxMsgTypeMeteringMsg
}

func (m *MeteringMsg) Bytes(txSerializer tx.TxSerializer) []byte {
	bytes, _ := txSerializer.MarshalJSON(m)

	return bytes
}

func (m *MeteringMsg) SetSecretKey(sk ankrcrypto.SecretKey) {

}

func (m *MeteringMsg) SecretKey() ankrcrypto.SecretKey {
	return &ankrcrypto.SecretKeyPem{}
}

func (m *MeteringMsg) PermitKey(store appstore.AppStore, pubKey []byte) bool {
	return true
}

func (m *MeteringMsg) ProcessTx(context tx.ContextTx, metric gas.GasMetric, isOnlyCheck bool) (uint32, string, []cmn.KVPair) {
	if isOnlyCheck {
		return code.CodeTypeOK, "", nil
	}

	context.AppStore().SetMetering(m.DCName, m.NSName, m.Value)

	context.AppStore().IncNonce(m.FromAddr)

	tvalue := time.Now().UnixNano()
	tags := []cmn.KVPair{
		{Key: []byte("app.fromaddress"), Value: []byte(m.FromAddr)},
		{Key: []byte("app.metering"), Value: []byte(m.DCName + ":" + m.NSName)},
		{Key: []byte("app.timestamp"), Value: []byte(strconv.FormatInt(tvalue, 10))},
		{Key: []byte("app.type"), Value: []byte(txcmm.TxMsgTypeMeteringMsg)},
	}

	return code.CodeTypeOK, "", tags
}


