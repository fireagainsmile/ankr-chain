package contract

import (
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/Ankr-network/ankr-chain/common/code"
	ankrcrypto "github.com/Ankr-network/ankr-chain/crypto"
	"github.com/Ankr-network/ankr-chain/store/appstore"
	"github.com/Ankr-network/ankr-chain/tx"
	txcmm "github.com/Ankr-network/ankr-chain/tx/common"
	ankrtypes "github.com/Ankr-network/ankr-chain/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

const (
	MIN_TOKEN_SEND = "5000000000000000000" // 5 tokens
)

type ContractDeployMsg struct {
	FromAddr string   `json:"fromaddr"`
	Name     string   `json:"name"`
	Codes    []byte   `json:"codes"`
	CodesDesc string  `json:"codesdesc"`
	gasUsed  *big.Int
}

type signContractDeployMsg struct {
	FromAddr  string   `json:"fromaddr"`
	Name      string   `json:"name"`
	Codes     []byte   `json:"codes"`
	CodesDesc string  `json:"codesdesc"`
}

func (sc signContractDeployMsg) bytes(txSerializer tx.TxSerializer) ([]byte, error) {
	return txSerializer.MarshalJSON(&sc)
}

func (cd *ContractDeployMsg) SignerAddr() []string {
	return []string {cd.FromAddr}
}

func (cd *ContractDeployMsg) GasWanted() int64 {
	return 0
}

func (cd *ContractDeployMsg) GasUsed() int64 {
	gasUsed, _ := strconv.ParseInt(MIN_TOKEN_SEND, 0, 64)

	return gasUsed
}

func (cd *ContractDeployMsg) Type() string {
	return txcmm.TxMsgTypeContractDeployMsg
}

func (cd *ContractDeployMsg) signMsg() *signContractDeployMsg {
	return &signContractDeployMsg{FromAddr: cd.FromAddr, Name: cd.Name, Codes: cd.Codes}
}

func (cd *ContractDeployMsg) Bytes(txSerializer tx.TxSerializer) []byte {
	bytes, _ :=  cd.signMsg().bytes(txSerializer)
	return bytes
}

func (cd *ContractDeployMsg) SetSecretKey(sk ankrcrypto.SecretKey) {

}

func (cd *ContractDeployMsg) SecretKey() ankrcrypto.SecretKey {
	return &ankrcrypto.SecretKeyEd25519{}
}

func (cd *ContractDeployMsg) PermitKey(store appstore.AppStore, pubKey []byte) bool {
	return true
}

func (cd *ContractDeployMsg) SpendGas(gas *big.Int) bool {
	cd.gasUsed.Add(cd.gasUsed, gas)
	return true
}

func (cd *ContractDeployMsg) SenderAddr() string {
	return cd.FromAddr
}

func (cd *ContractDeployMsg) ProcessTx(context tx.ContextTx, isOnlyCheck bool) (uint32, string, []cmn.KVPair){
	if len(cd.FromAddr) != ankrtypes.KeyAddressLen {
		return  code.CodeTypeInvalidAddress, fmt.Sprintf("ContractDeployMsg ProcessTx, unexpected from address. Got %s, addr len=%d", cd.FromAddr, len(cd.FromAddr)), nil
	}

	if len(cd.Codes) <= ankrtypes.CodePrefixLen {
		return code.CodeTypeContractInvalidCodeSize, fmt.Sprintf("ContractDeployMsg ProcessTx, invalid code size, Got %v, code size=%d", cd.Codes, len(cd.Codes)), nil
	}

	nonce, _ := context.AppStore().Nonce(cd.FromAddr)
	contractAddr := ankrcrypto.CreateContractAddress(cd.FromAddr, nonce)
    cInfo, err := context.AppStore().LoadContract(contractAddr)
    if err != nil {
    	return code.CodeTypeLoadContractErr, fmt.Sprintf("ContractDeployMsg ProcessTx, load contract err: contractAddr=%s", contractAddr), nil
	} else if cInfo != nil {
		return code.CodeTypeContractAddrTakenUp, fmt.Sprintf("ContractDeployMsg ProcessTx, the contract adress has been taken up:contractAddr=%s", contractAddr), nil
	}

    cInfo = &ankrtypes.ContractInfo{contractAddr, cd.Name, cd.FromAddr, cd.Codes, cd.CodesDesc}
    context.AppStore().SaveContract(contractAddr, cInfo)

	if isOnlyCheck {
		return code.CodeTypeOK, "", nil
	}

	context.AppStore().IncNonce(cd.FromAddr)

	tvalue := time.Now().UnixNano()
	tags := []cmn.KVPair{
		{Key: []byte("app.fromaddress"), Value: []byte(cd.FromAddr)},
		{Key: []byte("app.timestamp"), Value: []byte(strconv.FormatInt(tvalue, 10))},
		{Key: []byte("app.type"), Value: []byte(txcmm.TxMsgTypeContractDeployMsg)},
	}

	return code.CodeTypeOK, contractAddr, tags
}