package integratetesting

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Ankr-network/ankr-chain/account"
	"github.com/Ankr-network/ankr-chain/tool/cli/cmd"
	"github.com/Ankr-network/ankr-chain/tx/contract"
	"github.com/Ankr-network/ankr-chain/tx/metering"
	"io/ioutil"
	"math/big"
	"sync"
	"testing"

	"github.com/Ankr-network/ankr-chain/client"
	ankrcmm "github.com/Ankr-network/ankr-chain/common"
	"github.com/Ankr-network/ankr-chain/crypto"
	"github.com/Ankr-network/ankr-chain/tx/serializer"
	"github.com/Ankr-network/ankr-chain/tx/token"
	"github.com/stretchr/testify/assert"
)

const TEST_CERT = `
-----BEGIN CERTIFICATE-----
MIICKDCCAc6gAwIBAgIUVoDB7Av8NH3bhQDPajOX/AHq/zIwCgYIKoZIzj0EAwIw
dDELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAkNBMQswCQYDVQQHEwJTRjEUMBIGA1UE
CRMLTUlTU0lPTiBTVC4xDjAMBgNVBBETBTk0MTA1MQ4wDAYDVQQKEwVIVUJDQTEV
MBMGA1UEAxMMbXlodWItY2EuY29tMB4XDTE5MDYyNDA3NDk1NloXDTI5MDYyNDA3
NDk1NlowfTELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAkNBMQswCQYDVQQHEwJTRjEU
MBIGA1UECRMLTUlTU0lPTiBTVC4xDjAMBgNVBBETBTk0MTA1MRMwEQYDVQQKEwpE
YXRhQ2VudGVyMRkwFwYDVQQDExBteWRhdGFjZW50ZXIuY29tMFkwEwYHKoZIzj0C
AQYIKoZIzj0DAQcDQgAEE4x4SoWjyQit98+NDaAApQIbNIUOh/wGi4rR6EmcGmFa
qKW0jHxoYr3093CQHQ5X+BVVAjsLZCSy5melIcgPLqM1MDMwDgYDVR0PAQH/BAQD
AgeAMBMGA1UdJQQMMAoGCCsGAQUFBwMCMAwGA1UdEwEB/wQCMAAwCgYIKoZIzj0E
AwIDSAAwRQIhAPqre8XQqNr6JFvEhfaZz5XHf7854zDC4H/wmLcRv5b3AiAGgiuI
PvDQFLYt8PkvJk9hH2ynYEyI6zId1KFGxBrd/g==
-----END CERTIFICATE-----`

const TEST_KEY = `
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIHMyEp01U2qMTNYLdQEyc9NX8F9JowMM7ODVD9ap77ENoAoGCCqGSM49
AwEHoUQDQgAEE4x4SoWjyQit98+NDaAApQIbNIUOh/wGi4rR6EmcGmFaqKW0jHxo
Yr3093CQHQ5X+BVVAjsLZCSy5melIcgPLg==
-----END EC PRIVATE KEY-----`

var(
	testChainId = ankrcmm.ChainID("test-chain-Aavjfy")
	localUrl = "localhost:26657"
	remoteUrl = ""
	adminAddress = "B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67"
	adminPriv = "wmyZZoMedWlsPUDVCOy+TiVcrIBPcn3WJN8k5cPQgIvC8cbcR10FtdAdzIlqXQJL9hBw1i0RsVjF6Oep/06Ezg=="
	txSerializer = serializer.NewTxSerializerCDC()
)
var defaultTxMsgHeader = client.TxMsgHeader{
	ChID: testChainId,
	GasLimit: new(big.Int).SetUint64(1000).Bytes(),
	GasPrice: ankrcmm.Amount{ankrcmm.Currency{"ANKR", 18}, new(big.Int).SetUint64(100000000000000000).Bytes()},
	Memo: "test transfer",
	Version: "1.0",
}

// base 1000000000000000000
// Procedure:
// 1. init one account with coins
// 2. send transaction from the init account to the receive account
// 3. check the receive account balance
func TestSendToOneRepeatedly(t *testing.T)  {
	c := client.NewClient(localUrl)
	msgHeader := defaultTxMsgHeader
	_, toAddr := cmd.GenAccount()
	// init one account with balance
	initPriv, initAddress, err := initAccountWithAmout("600000000000000000000")
	if err != nil {
		t.Error(err)
	}
	// send transaction to one account repeatedly
	initKey := crypto.NewSecretKeyEd25519(initPriv)
	repeatedTimes := 3
	repeatAmount, _ := new(big.Int).SetString("10000000000000000000", 10) //10 ankr token
	repeatMsg := newTransferMsg(initAddress, toAddr, repeatAmount)
	for i := 0; i < repeatedTimes ; i++ {
		repeatBuilder := client.NewTxMsgBuilder(msgHeader, repeatMsg, txSerializer, initKey)
		_, _, _, err = repeatBuilder.BuildAndCommit(c)
		if err != nil {
			t.Log(err)
			return
		}
	}
	expectAmount := new(big.Int).Mul(new(big.Int).SetUint64(uint64(repeatedTimes)), repeatAmount)
	err = checkBlance(c, toAddr, expectAmount.String())
	if err != nil {
		t.Error(err)
	}
}

type Account struct {
	address string
	privateKey string
}

// Procedure:
// 1. init multiple account with coins
// 2. send transaction from the init accounts to one receive account simultaneously
// 3. check the receive account balance
func TestMultipleAccSendToOne(t *testing.T)  {
	var accounts []Account
	accountNum := 5
	for i := 0; i < accountNum; i++ {
		priv, addr ,err := initAccountWithAmout("10000000000000000000") //10 ankr token
		if err != nil {
			t.Error(err)
		}
		accounts = append(accounts, Account{address:addr, privateKey: priv})
	}

	_, toAddress := cmd.GenAccount()
	msgHeader := defaultTxMsgHeader
	sendAmount, _ := new(big.Int).SetString("7000000000000000000", 10)
	var wg sync.WaitGroup
	wg.Add(accountNum)
	for i := 0 ; i < accountNum ; i++ {
		go func(index int) {
			acc := accounts[index]
			c := client.NewClient(localUrl)
			txMsg := newTransferMsg(acc.address, toAddress, sendAmount )
			key := crypto.NewSecretKeyEd25519(acc.privateKey)
			builder := client.NewTxMsgBuilder(msgHeader, txMsg, txSerializer, key)
			_, _, _, err := builder.BuildAndCommit(c)
			wg.Done()
			if err != nil {
				t.Error(err)
			}
		}(i)
	}
	wg.Wait()
	cl := client.NewClient(localUrl)
	expectAmount := new(big.Int).Mul(sendAmount, new(big.Int).SetUint64(uint64(accountNum)))
	err := checkBlance(cl, toAddress, expectAmount.String())
	if err != nil {
		t.Error(err)
	}
}

func TestTxWithNonce(t *testing.T)  {
	c := client.NewClient("localhost:26657")
	msgHeader := defaultTxMsgHeader
	initPriv, initAddr, err := initAccountWithAmout("30000000000000000000") // 30 ankr token
	if err != nil {
		t.Error(err)
	}
	_, toAddr := cmd.GenAccount()
	sendAmount, _ := new(big.Int).SetString("2000000000000000000", 10) //2 ankr token
	txMsg := newTransferMsg(initAddr, toAddr, sendAmount)
	key := crypto.NewSecretKeyEd25519(initPriv)
	builder := client.NewTxMsgBuilder(msgHeader, txMsg, txSerializer, key)
	nonceReq := new(ankrcmm.NonceQueryReq)
	nonceReq.Address = initAddr
	nonceResp := new(ankrcmm.NonceQueryResp)
	err = c.Query("/store/nonce", nonceReq, nonceResp)
	if err != nil {
		t.Error(err)
	}
	nonce := nonceResp.Nonce
	var txs [][]byte
	for i := 0; i < 4; i++ {
		tx, err := builder.BuildOnly(nonce + uint64(i))
		if err != nil {
			t.Error(err)
		}
		txs = append(txs, tx)
	}

	//skip nonce
	_, _, _, err = c.BroadcastTxCommit(txs[3])
	if err == nil {
		t.Error(errors.New("expecting error"))
	}
	t.Log(err)
	// duplicate nonce
	_, _, _, err = c.BroadcastTxCommit(txs[0])
	if err != nil {
		t.Error(err)
	}
	_, _, _, err = c.BroadcastTxCommit(txs[0])
	if err == nil {
		t.Error(errors.New("expecting error"))
	}
	for i := 0; i < 3 ; i ++ {
		_, _, _, err = c.BroadcastTxCommit(txs[i])
		if err != nil {
			t.Log(err)
		}
	}
	err = checkBlance(c, toAddr, "")
	t.Log(err)
}

func TestTxTransferWithNode(t *testing.T) {
	c := client.NewClient("localhost:26657")
	msgHeader := client.TxMsgHeader{
		ChID: testChainId,
		GasLimit: new(big.Int).SetUint64(1000).Bytes(),
		GasPrice: ankrcmm.Amount{ankrcmm.Currency{"ANKR", 18}, new(big.Int).SetUint64(100000000000000000).Bytes()},
		Memo: "test transfer",
		Version: "1.0",
	}

	tfMsg := &token.TransferMsg{FromAddr: "B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67",
		ToAddr:  "065E37B3FC243B9FABB1519AB876E7632C510DC9324031",
		Amounts: []ankrcmm.Amount{ankrcmm.Amount{ankrcmm.Currency{"ANKR", 18}, new(big.Int).SetUint64(6000000000000000000).Bytes()}},
	}

	txSerializer := serializer.NewTxSerializerCDC()

	key := crypto.NewSecretKeyEd25519("wmyZZoMedWlsPUDVCOy+TiVcrIBPcn3WJN8k5cPQgIvC8cbcR10FtdAdzIlqXQJL9hBw1i0RsVjF6Oep/06Ezg==")

	builder := client.NewTxMsgBuilder(msgHeader, tfMsg,  txSerializer, key)

	txHash, cHeight, _, err := builder.BuildAndCommit(c)

	assert.Equal(t, err, nil)

	t.Logf("TestTxTransferWithNode sucessful: txHash=%s, cHeight=%d", txHash, cHeight)

	resp := &ankrcmm.BalanceQueryResp{}
	c.Query("/store/balance", &ankrcmm.BalanceQueryReq{"065E37B3FC243B9FABB1519AB876E7632C510DC9324031", "ANKR"}, resp)

	t.Logf("addr=065E37B3FC243B9FABB1519AB876E7632C510DC9324031, bal=%s", resp.Amount)

	resp = &ankrcmm.BalanceQueryResp{}
	c.Query("/store/balance", &ankrcmm.BalanceQueryReq{"B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67", "ANKR"}, resp)

	t.Logf("addr=B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67, bal=%s", resp.Amount)
}

func TestCertMsgWithNode(t *testing.T) {
	c := client.NewClient("localhost:26657")

	msgHeader := client.TxMsgHeader{
		ChID: "test-chain-Aavjfy",
		GasLimit: new(big.Int).SetUint64(1000).Bytes(),
		GasPrice: ankrcmm.Amount{ankrcmm.Currency{"ANKR", 18}, new(big.Int).SetUint64(10000000000000000).Bytes()},
		Memo: "test CertMsg",
		Version: "1.0",
	}

	pubBS64 := account.AccountManagerInstance().AdminOpAccount(ankrcmm.AccountAdminMetering)
	addrFrom, err := ankrcmm.AddressByPublicKey(pubBS64)

	t.Logf("certMsgFromAddr=%s", addrFrom)
	assert.Equal(t, err, nil)

	certMsg := &metering.SetCertMsg{FromAddr: addrFrom,
		DCName: "dc1",
	    PemBase64: TEST_CERT,
	}

	txSerializer := serializer.NewTxSerializerCDC()

	key := crypto.NewSecretKeyEd25519("wmyZZoMedWlsPUDVCOy+TiVcrIBPcn3WJN8k5cPQgIvC8cbcR10FtdAdzIlqXQJL9hBw1i0RsVjF6Oep/06Ezg==")

	builder := client.NewTxMsgBuilder(msgHeader, certMsg,  txSerializer, key)

	txHash, cHeight, _, err := builder.BuildAndCommit(c)

	assert.Equal(t, err, nil)

	t.Logf("TestCertMsgWithNode:94 sucessful: txHash=%s, cHeight=%d", txHash, cHeight)

	resp := &ankrcmm.CertKeyQueryResp{}
	c.Query("/store/certkey", &ankrcmm.CertKeyQueryReq{"dc1"}, resp)

	t.Logf("pembase64=%s", resp.PEMBase64)
}

func TestMeteringWithNode(t *testing.T) {
	c := client.NewClient("localhost:26657")

	msgHeader := client.TxMsgHeader{
		ChID: "test-chain-Aavjfy",
		GasLimit: new(big.Int).SetUint64(1000).Bytes(),
		GasPrice: ankrcmm.Amount{ankrcmm.Currency{"ANKR", 18}, new(big.Int).SetUint64(10000000000000).Bytes()},
		Memo: "test metering",
		Version: "1.0",
	}

	resp := &ankrcmm.CertKeyQueryResp{}
	c.Query("/store/certkey", &ankrcmm.CertKeyQueryReq{"dc1"}, resp)

	key := crypto.NewSecretKeyPem(TEST_KEY, resp.PEMBase64,"@mert:"+"dc1_"+"ns1")

	addr, _ := key.Address()

	t.Logf("meteringtMsgFromAddr=%s", string(addr))

	certMsg := &metering.MeteringMsg{FromAddr: string(addr),
		DCName: "dc1",
		NSName: "ns1",
		Value: "value1",
	}

	txSerializer := serializer.NewTxSerializerCDC()

	builder := client.NewTxMsgBuilder(msgHeader, certMsg,  txSerializer, key)

	txHash, cHeight, _, err := builder.BuildAndCommit(c)

	assert.Equal(t, err, nil)

	t.Logf("TestCertMsgWithNode:94 sucessful: txHash=%s, cHeight=%d", txHash, cHeight)

	respMetering := &ankrcmm.MeteringQueryResp{}
	c.Query("/store/metering", &ankrcmm.MeteringQueryReq{"dc1", "ns1"}, respMetering)

	t.Logf("value=%s", respMetering.Value)
}

func TestContractDeployWithNode(t *testing.T) {
	c := client.NewClient("localhost:26657")

	msgHeader := client.TxMsgHeader{
		ChID: "test-chain-Aavjfy",
		GasLimit: new(big.Int).SetUint64(5000000).Bytes(),
		GasPrice: ankrcmm.Amount{ankrcmm.Currency{"ANKR", 18}, new(big.Int).SetUint64(10000000000000000).Bytes()},
		Memo: "test ContractDeploy",
		Version: "1.0",
	}

	rawBytes, err := ioutil.ReadFile("TestContract.wasm")
	if err != nil {
		t.Errorf("can't read wasm file: %s", err.Error())
	}

	cdMsg := &contract.ContractDeployMsg{FromAddr: "B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67",
		Name:     "TestContract",
		Codes:     rawBytes,
	    CodesDesc: "",
	}

	txSerializer := serializer.NewTxSerializerCDC()

	key := crypto.NewSecretKeyEd25519("wmyZZoMedWlsPUDVCOy+TiVcrIBPcn3WJN8k5cPQgIvC8cbcR10FtdAdzIlqXQJL9hBw1i0RsVjF6Oep/06Ezg==")

	builder := client.NewTxMsgBuilder(msgHeader, cdMsg,  txSerializer, key)

	txHash, cHeight, contractAddr, err := builder.BuildAndCommit(c)

	assert.Equal(t, err, nil)

	t.Logf("TestTxTransferWithNode sucessful: txHash=%s, cHeight=%d, contractAddr=%s", txHash, cHeight, contractAddr)

	resp := &ankrcmm.ContractQueryResp{}
	c.Query("/store/contract", &ankrcmm.ContractQueryReq{contractAddr}, resp)

	t.Logf("conract=%v", resp)
}

// Procedure:
// 1. read smart contract binary from file
// 2. copy part of contract binary into tranasction message
// 3. send transaction and check response
// 4. check validator status
func TestDeployBadContract(t *testing.T) {
	c := client.NewClient("localhost:26657")

	msgHeader := client.TxMsgHeader{
		ChID: "test-chain-ZFaVSj",
		GasLimit: new(big.Int).SetUint64(100000000000).Bytes(),
		GasPrice: ankrcmm.Amount{ankrcmm.Currency{"ANKR", 18}, new(big.Int).SetUint64(1000000000000000000).Bytes()},
		Memo: "test ContractDeploy",
		Version: "1.0",
	}

	rawBytes, err := ioutil.ReadFile("TestContract.wasm")
	if err != nil {
		t.Errorf("can't read wasm file: %s", err.Error())
	}

	// only copy part of contract into transaction msg
	length := len(rawBytes)
	rawBytes = rawBytes[:length - 8]

	cdMsg := &contract.ContractDeployMsg{FromAddr: "B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67",
		Name:     "TestContract",
		Codes:     rawBytes,
		CodesDesc: "",
	}

	txSerializer := serializer.NewTxSerializerCDC()

	key := crypto.NewSecretKeyEd25519("wmyZZoMedWlsPUDVCOy+TiVcrIBPcn3WJN8k5cPQgIvC8cbcR10FtdAdzIlqXQJL9hBw1i0RsVjF6Oep/06Ezg==")

	builder := client.NewTxMsgBuilder(msgHeader, cdMsg,  txSerializer, key)

	txHash, cHeight, contractAddr, err := builder.BuildAndCommit(c)

	assert.Equal(t, err, nil)

	t.Logf("TestTxTransferWithNode sucessful: txHash=%s, cHeight=%d, contractAddr=%s", txHash, cHeight, contractAddr)

	resp := &ankrcmm.ContractQueryResp{}
	c.Query("/store/contract", &ankrcmm.ContractQueryReq{contractAddr}, resp)

	t.Logf("conract=%v", resp)
}

func TestContractInvokeWithNode(t *testing.T) {
	c := client.NewClient("localhost:26657")

	msgHeader := client.TxMsgHeader{
		ChID: "test-chain-xPPj8k",
		GasLimit: new(big.Int).SetUint64(1000).Bytes(),
		GasPrice: ankrcmm.Amount{ankrcmm.Currency{"ANKR", 18}, new(big.Int).SetUint64(10000000000000).Bytes()},
		Memo: "test ContractInvoke",
		Version: "1.0",
	}

	jsonArg := "[{\"index\":1,\"Name\":\"args\",\"ParamType\":\"string\",\"Value\":{\"testStr\":\"testFuncWithInt arg\"}}]"

	cdMsg := &contract.ContractInvokeMsg{
		FromAddr: "B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67",
		ContractAddr: "AD11BED29F81AE1DD51DE7127A5A99859DB60E1E7B19B2",
		Method:       "testFuncWithString",
		Args:         jsonArg,
		RtnType:      "string",
	}

	txSerializer := serializer.NewTxSerializerCDC()

	key := crypto.NewSecretKeyEd25519("wmyZZoMedWlsPUDVCOy+TiVcrIBPcn3WJN8k5cPQgIvC8cbcR10FtdAdzIlqXQJL9hBw1i0RsVjF6Oep/06Ezg==")

	builder := client.NewTxMsgBuilder(msgHeader, cdMsg,  txSerializer, key)

	txHash, cHeight, contractResultJson, err := builder.BuildAndCommit(c)

	assert.Equal(t, err, nil)

	var contractR ankrcmm.ContractResult
	json.Unmarshal([]byte(contractResultJson), &contractR)

	t.Logf("TestTxTransferWithNode sucessful: txHash=%s, cHeight=%d, contractR=%v", txHash, cHeight, contractR)
}

func TestContractDeployWithNodePattern2(t *testing.T) {
	c := client.NewClient("localhost:26657")

	msgHeader := client.TxMsgHeader{
		ChID: "test-chain-b4lpJu",
		GasLimit: new(big.Int).SetUint64(1000).Bytes(),
		GasPrice: ankrcmm.Amount{ankrcmm.Currency{"ANKR", 18}, new(big.Int).SetUint64(10000000000000).Bytes()},
		Memo: "test ContractDeploy",
		Version: "1.0",
	}

	//rawBytes, err := ioutil.ReadFile("F:/GoPath/src/github.com/Ankr-network/ankr-chain/contract/example/cpp/TestContract2.wasm")
	rawBytes, err := ioutil.ReadFile("F:/GoPath/src/github.com/Ankr-network/ankr-chain/contract/example/cpp/ERC20.wasm")
	if err != nil {
		t.Errorf("can't read wasm file: %s", err.Error())
	}

	cdMsg := &contract.ContractDeployMsg{FromAddr: "B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67",
		Name:     "TestContract",
		Codes:     rawBytes,
		CodesDesc: "",
	}

	txSerializer := serializer.NewTxSerializerCDC()

	key := crypto.NewSecretKeyEd25519("wmyZZoMedWlsPUDVCOy+TiVcrIBPcn3WJN8k5cPQgIvC8cbcR10FtdAdzIlqXQJL9hBw1i0RsVjF6Oep/06Ezg==")

	builder := client.NewTxMsgBuilder(msgHeader, cdMsg,  txSerializer, key)

	txHash, cHeight, contractAddr, err := builder.BuildAndCommit(c)

	assert.Equal(t, err, nil)

	t.Logf("TestTxTransferWithNode sucessful: txHash=%s, cHeight=%d, contractAddr=%s", txHash, cHeight, contractAddr)

	resp := &ankrcmm.ContractQueryResp{}
	c.Query("/store/contract", &ankrcmm.ContractQueryReq{contractAddr}, resp)

	t.Logf("conract=%v", resp)
}

func TestContractInvokeWithNodePattern2(t *testing.T) {
	c := client.NewClient("localhost:26657")

	msgHeader := client.TxMsgHeader{
		ChID: "test-chain-qJeO7E",
		GasLimit: new(big.Int).SetUint64(1000).Bytes(),
		GasPrice: ankrcmm.Amount{ankrcmm.Currency{"ANKR", 18}, new(big.Int).SetUint64(10000000000000).Bytes()},
		Memo: "test ContractInvoke",
		Version: "1.0",
	}

	jsonArg := "[{\"index\":1,\"Name\":\"args\",\"ParamType\":\"string\",\"Value\":\"testFuncWithInt arg\"}]"

	cdMsg := &contract.ContractInvokeMsg{
		FromAddr: "B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67",
		ContractAddr: "A277D0BD075656C3DBF92F9FEDC396AFFC75C95B9CF6D6",
		Method:       "testFuncWithString",
		Args:         jsonArg,
		RtnType:      "string",
	}

	txSerializer := serializer.NewTxSerializerCDC()

	key := crypto.NewSecretKeyEd25519("wmyZZoMedWlsPUDVCOy+TiVcrIBPcn3WJN8k5cPQgIvC8cbcR10FtdAdzIlqXQJL9hBw1i0RsVjF6Oep/06Ezg==")

	builder := client.NewTxMsgBuilder(msgHeader, cdMsg,  txSerializer, key)

	txHash, cHeight, contractResultJson, err := builder.BuildAndCommit(c)

	assert.Equal(t, err, nil)

	var contractR ankrcmm.ContractResult
	json.Unmarshal([]byte(contractResultJson), &contractR)

	t.Logf("TestTxTransferWithNode sucessful: txHash=%s, cHeight=%d, contractR=%v", txHash, cHeight, contractR)
}

func TestQueryAccountInfoWithNode(t *testing.T) {
	c := client.NewClient("localhost:26657")

	resp := &ankrcmm.AccountQueryResp{}
	c.Query("/store/account", &ankrcmm.AccountQueryReq{"B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67"}, resp)

	t.Logf("account=%v", resp)
}


// helper functions used in test case
func newTransferMsg(from, to string, amount *big.Int) *token.TransferMsg {
	return &token.TransferMsg{
		FromAddr: from,
		ToAddr:  to,
		Amounts: []ankrcmm.Amount{ankrcmm.Amount{ankrcmm.Currency{"ANKR", 18}, amount.Bytes()}},
	}
}

func checkBlance(cl *client.Client, address string, expect string) error {
	balReq := new(ankrcmm.BalanceQueryReq)
	balReq.Address = address
	balReq.Symbol = "ANKR"
	balResp := new(ankrcmm.BalanceQueryResp)
	err := cl.Query("/store/balance", balReq, balResp)
	if err != nil {
		return err
	}
	if balResp.Amount != expect {
		return errors.New(fmt.Sprintf("expect %s, got %s", expect, balResp.Amount))
	}
	return nil
}

//generate account with coins
func initAccountWithAmout(amount string) (priv, address string,err error)  {
	c := client.NewClient(localUrl)
	msgHeader := defaultTxMsgHeader
	priv, address = cmd.GenAccount()
	initAmount, _ := new(big.Int).SetString(amount, 10)
	initMsg := newTransferMsg(adminAddress, address, initAmount)
	adminKey := crypto.NewSecretKeyEd25519(adminPriv)
	//transaction builder
	builder := client.NewTxMsgBuilder(msgHeader, initMsg,  txSerializer, adminKey)
	_, _, _, err = builder.BuildAndCommit(c)
	if err != nil {
		return "", "", err
	}
	//check balance
	err = checkBlance(c, address, amount)
	if err != nil {
		return "", "", err
	}
	return
}