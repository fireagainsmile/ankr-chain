package integratetesting

import (
	"github.com/Ankr-network/ankr-chain/account"
	"github.com/Ankr-network/ankr-chain/tx/metering"
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

func TestTxTransferWithNode(t *testing.T) {
	c := client.NewClient("localhost:26657")

	msgHeader := client.TxMsgHeader{
		ChID: "test-chain-hQYhLJ",
		GasLimit: new(big.Int).SetUint64(1000).Bytes(),
		GasPrice: ankrcmm.Amount{ankrcmm.Currency{"ANKR", 18}, new(big.Int).SetUint64(10000000000000).Bytes()},
		Memo: "test transfer",
		Version: "1.0",
	}

	tfMsg := &token.TransferMsg{FromAddr: "B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67",
		ToAddr:  "454D92DC842F532683E820DF6C3784473AD9CCF222D8FB",
		Amounts: []ankrcmm.Amount{ankrcmm.Amount{ankrcmm.Currency{"ANKR", 18}, new(big.Int).SetUint64(6000000000000000000).Bytes()}},
	}

	txSerializer := serializer.NewTxSerializerCDC()

	key := crypto.NewSecretKeyEd25519("wmyZZoMedWlsPUDVCOy+TiVcrIBPcn3WJN8k5cPQgIvC8cbcR10FtdAdzIlqXQJL9hBw1i0RsVjF6Oep/06Ezg==")

	builder := client.NewTxMsgBuilder(msgHeader, tfMsg,  txSerializer, key)

	txHash, cHeight, err := builder.BuildAndCommit(c)

	assert.Equal(t, err, nil)

	t.Logf("TestTxTransferWithNode sucessful: txHash=%s, cHeight=%d", txHash, cHeight)

	resp := &ankrcmm.BalanceQueryResp{}
	c.Query("/store/balance", &ankrcmm.BalanceQueryReq{"B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67", "ANKR"}, resp)

	t.Logf("bal=%s", resp.Amount)
}

func TestCertMsgWithNode(t *testing.T) {
	c := client.NewClient("localhost:26657")

	msgHeader := client.TxMsgHeader{
		ChID: "test-chain-hQYhLJ",
		GasLimit: new(big.Int).SetUint64(1000).Bytes(),
		GasPrice: ankrcmm.Amount{ankrcmm.Currency{"ANKR", 18}, new(big.Int).SetUint64(10000000000000).Bytes()},
		Memo: "test transfer",
		Version: "1.0",
	}

	pubBS64 := account.AccountManagerInstance().AdminOpAccount(ankrcmm.AccountAdminMetering)
	addrFrom, err := ankrcmm.AddressByPublicKey(pubBS64)

	t.Logf("certMsgFromAddr=%s", addrFrom)
	assert.Equal(t, err, nil)

	certMsg := &metering.SetCertMsg{FromAddr: addrFrom,
		DCName: "dc1",
		NSName: "ns1",
	    PemBase64: TEST_CERT,
	}

	txSerializer := serializer.NewTxSerializerCDC()

	key := crypto.NewSecretKeyEd25519("wmyZZoMedWlsPUDVCOy+TiVcrIBPcn3WJN8k5cPQgIvC8cbcR10FtdAdzIlqXQJL9hBw1i0RsVjF6Oep/06Ezg==")

	builder := client.NewTxMsgBuilder(msgHeader, certMsg,  txSerializer, key)

	txHash, cHeight, err := builder.BuildAndCommit(c)

	assert.Equal(t, err, nil)

	t.Logf("TestCertMsgWithNode:94 sucessful: txHash=%s, cHeight=%d", txHash, cHeight)

	resp := &ankrcmm.CertKeyQueryResp{}
	c.Query("/store/certkey", &ankrcmm.CertKeyQueryReq{"dc1", "ns1"}, resp)

	t.Logf("pembase64=%s", resp.PEMBase64)
}
