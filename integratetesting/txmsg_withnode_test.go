package integratetesting

import (
	"encoding/json"
	"github.com/Ankr-network/ankr-chain/account"
	"github.com/Ankr-network/ankr-chain/tx/contract"
	"github.com/Ankr-network/ankr-chain/tx/metering"
	"io/ioutil"
	"math/big"
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

func TestTxTransferWithNode(t *testing.T) {
	c := client.NewClient("localhost:26657")

	msgHeader := client.TxMsgHeader{
		ChID: "test-chain-b4lpJu",
		GasLimit: new(big.Int).SetUint64(1000).Bytes(),
		GasPrice: ankrcmm.Amount{ankrcmm.Currency{"ANKR", 18}, new(big.Int).SetUint64(10000000000000).Bytes()},
		Memo: "test transfer",
		Version: "1.0",
	}

	tfMsg := &token.TransferMsg{FromAddr: "B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67",
		ToAddr:  "065E37B3FC243B9FABB1519AB876E7632C510DC9324031",
		Amounts: []ankrcmm.Amount{ankrcmm.Amount{ankrcmm.Currency{"TESTCOIN", 18}, new(big.Int).SetUint64(6000000000000000000).Bytes()}},
	}

	txSerializer := serializer.NewTxSerializerCDC()

	key := crypto.NewSecretKeyEd25519("wmyZZoMedWlsPUDVCOy+TiVcrIBPcn3WJN8k5cPQgIvC8cbcR10FtdAdzIlqXQJL9hBw1i0RsVjF6Oep/06Ezg==")

	builder := client.NewTxMsgBuilder(msgHeader, tfMsg,  txSerializer, key)

	txHash, cHeight, _, err := builder.BuildAndCommit(c)

	assert.Equal(t, err, nil)

	t.Logf("TestTxTransferWithNode sucessful: txHash=%s, cHeight=%d", txHash, cHeight)

	resp := &ankrcmm.BalanceQueryResp{}
	c.Query("/store/balance", &ankrcmm.BalanceQueryReq{"065E37B3FC243B9FABB1519AB876E7632C510DC9324031", "TESTCOIN"}, resp)

	t.Logf("bal=%s", resp.Amount)
}

func TestCertMsgWithNode(t *testing.T) {
	c := client.NewClient("localhost:26657")

	msgHeader := client.TxMsgHeader{
		ChID: "test-chain-hQYhLJ",
		GasLimit: new(big.Int).SetUint64(1000).Bytes(),
		GasPrice: ankrcmm.Amount{ankrcmm.Currency{"ANKR", 18}, new(big.Int).SetUint64(10000000000000).Bytes()},
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
		ChID: "test-chain-hQYhLJ",
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
		ChID: "test-chain-b4lpJu",
		GasLimit: new(big.Int).SetUint64(1000).Bytes(),
		GasPrice: ankrcmm.Amount{ankrcmm.Currency{"ANKR", 18}, new(big.Int).SetUint64(10000000000000).Bytes()},
		Memo: "test ContractDeploy",
		Version: "1.0",
	}

	rawBytes, err := ioutil.ReadFile("F:/GoPath/src/github.com/Ankr-network/ankr-chain/contract/example/cpp/TestContract.wasm")
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
