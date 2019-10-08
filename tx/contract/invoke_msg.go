package contract

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/Ankr-network/ankr-chain/common/code"
	ankrcontext "github.com/Ankr-network/ankr-chain/context"
	"github.com/Ankr-network/ankr-chain/contract"
	ankrcrypto "github.com/Ankr-network/ankr-chain/crypto"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	"github.com/Ankr-network/ankr-chain/tx"
	txcmm "github.com/Ankr-network/ankr-chain/tx/common"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

type ContractInvokeMsg struct {
	FromAddr     string  `json:"fromaddr"`
	ContractAddr string  `json:"contractaddr"`
	Method       string  `json:"name"`
	Args         string  `json:"args"`
	RtnType      string  `json:"rtnType"`
	gasUsed      *big.Int
}

type signContractInvokeMsg struct {
	FromAddr     string  `json:"fromaddr"`
	ContractAddr string  `json:"contractaddr"`
	Method       string  `json:"name"`
	Args         string  `json:"args"`
	RtnType      string  `json:"rtnType"`
}

func (sc signContractInvokeMsg) bytes(txSerializer tx.TxSerializer) ([]byte, error) {
	return txSerializer.MarshalJSON(&sc)
}

func (cd *ContractInvokeMsg) SignerAddr() []string {
	return []string {cd.FromAddr}
}

func (cd *ContractInvokeMsg) GasWanted() int64 {
	return 0
}

func (cd *ContractInvokeMsg) GasUsed() int64 {
	gasUsed, _ := strconv.ParseInt(MIN_TOKEN_SEND, 0, 64)

	return gasUsed
}

func (cd *ContractInvokeMsg) Type() string {
	return txcmm.TxMsgTypeContractInvokeMsg
}

func (ci *ContractInvokeMsg) signMsg() *signContractInvokeMsg {
	return &signContractInvokeMsg{FromAddr: ci.FromAddr, ContractAddr: ci.ContractAddr, Method: ci.Method, Args: ci.Args}
}

func (ci *ContractInvokeMsg) Bytes(txSerializer tx.TxSerializer) []byte {
	bytes, _ :=  ci.signMsg().bytes(txSerializer)
	return bytes
}

func (ci *ContractInvokeMsg) SetSecretKey(sk ankrcrypto.SecretKey) {

}

func (ci *ContractInvokeMsg) SecretKey() ankrcrypto.SecretKey {
	return &ankrcrypto.SecretKeyEd25519{}
}

func (ci *ContractInvokeMsg) PermitKey(store appstore.AppStore, pubKey []byte) bool {
	return true
}

func (ci *ContractInvokeMsg) SpendGas(gas *big.Int) bool {
	ci.gasUsed.Add(ci.gasUsed, gas)
	return true
}

func (ci *ContractInvokeMsg) SenderAddr() string {
	return ci.FromAddr
}

func (ci *ContractInvokeMsg) ProcessTx(context tx.ContextTx, isOnlyCheck bool) (uint32, string, []cmn.KVPair) {
	if len(ci.FromAddr) != ankrtypes.KeyAddressLen {
		return code.CodeTypeInvalidAddress, fmt.Sprintf("ContractInvokeMsg ProcessTx, unexpected from address. Got %s, addr len=%d", ci.FromAddr, len(ci.FromAddr)), nil
	}

	if len(ci.ContractAddr) != ankrtypes.KeyAddressLen {
		return code.CodeTypeContractInvalidAddr, fmt.Sprintf("ContractInvokeMsg ProcessTx, unexpected contract address. Got %s, addr len=%d", ci.ContractAddr, len(ci.ContractAddr)), nil
	}

	cInfo, err := context.AppStore().LoadContract(ci.ContractAddr)
	if err != nil {
		return code.CodeTypeLoadContractErr, fmt.Sprintf("ContractInvokeMsg ProcessTx, load contract err: contractAddr=%s", ci.ContractAddr), nil
	} else if cInfo == nil {
		return code.CodeTypeContractAddrTakenUp, fmt.Sprintf("ContractInvokeMsg ProcessTx, the contract adress has been taken up:contractAddr=%s", ci.ContractAddr), nil
	}

	if isOnlyCheck {
		return code.CodeTypeOK, "", nil
	}

	var params []*ankrtypes.Param
	json.Unmarshal([]byte(ci.Args), &params)

	contractType    := ankrtypes.ContractType(cInfo.Codes[0])
	contractContext := ankrcontext.NewContextContract(ci, ci, context.AppStore())
	rtn, err := contract.Call(contractContext, contractType, cInfo.Codes[ankrtypes.CodePrefixLen:], cInfo.Name, ci.Method, params, ci.RtnType)
	if err != nil {
		return code.CodeTypeCallContractErr, fmt.Sprintf("call contract err: contract=%s, method=%s, err=%v", ci.ContractAddr, ci.Method, err), nil
	}

	if !rtn.IsSuccess {
		return code.CodeTypeCallContractErr, fmt.Sprintf("call contract err: contract=%s, method=%s", ci.ContractAddr, ci.Method), nil
	}

	context.AppStore().IncNonce(ci.FromAddr)

	tvalue := time.Now().UnixNano()
	tags := []cmn.KVPair{
		{Key: []byte("app.fromaddress"), Value: []byte(ci.FromAddr)},
		{Key: []byte("app.contractaddr"), Value: []byte(ci.ContractAddr)},
		{Key: []byte("app.method"), Value: []byte(ci.Method)},
		{Key: []byte("app.timestamp"), Value: []byte(strconv.FormatInt(tvalue, 10))},
		{Key: []byte("app.type"), Value: []byte(txcmm.TxMsgTypeContractInvokeMsg)},
	}

	rtnJson, _ := json.Marshal(rtn)

	return code.CodeTypeOK, string(rtnJson), tags
}
