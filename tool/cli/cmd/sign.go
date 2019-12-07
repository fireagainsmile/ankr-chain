package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/Ankr-network/ankr-chain/client"
	"github.com/Ankr-network/ankr-chain/crypto"
	"github.com/Ankr-network/ankr-chain/tx/serializer"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var(
	signFile = "signFile"
	signKeyStore = "signKeyStore"
)
// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "sign a raw transaction",
	Run: runSignTx,
}

func init() {
	err := addStringFlag(signCmd, signFile, fileParam, "f", "","file name of the json format transaction", required)
	if err != nil{
		panic(err)
	}

	err = addStringFlag(signCmd, signKeyStore, keystoreParam, "k", "","keystore file used to sign this transaction", required)
	if err != nil{
		panic(err)
	}
}

func runSignTx(cmd *cobra.Command, args []string)  {
	txRawFile := viper.GetString(signFile)
	keyFile := viper.GetString(signKeyStore)
	txBytes, err := ioutil.ReadFile(txRawFile)
	if err != nil {
		fmt.Println("Failed to read transaction info:", err.Error())
		return
	}
	keystorePath, err := getKeystoreFile(keyFile)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	privateKey := decryptPrivatekey(keystorePath)
	if privateKey == "" {
		fmt.Println("Error: Wrong keystore or password!")
		return
	}

	key := crypto.NewSecretKeyEd25519(privateKey)
	keyAddr, _ := key.Address()
	fromAddr := fmt.Sprintf("%X", keyAddr)

	var rawTx RawTransaction
	err = json.Unmarshal(txBytes, &rawTx)
	if err != nil {
		fmt.Println("Failed to unmarshal json file:", err.Error())
		return
	}

	msg := rawTx.TxMsg
	if msg.FromAddr != fromAddr {
		fmt.Println("can not match from address with the keystore")
		fmt.Println("Transaction from address:", msg.FromAddr)
		fmt.Println("Keystore address:", fromAddr)
		return
	}

	builder := client.NewTxMsgBuilder(*rawTx.Header, rawTx.TxMsg, serializer.NewTxSerializerCDC(), key)
	signedTxByte, err := builder.BuildOnly(rawTx.Nonce)
	if err != nil {
		fmt.Println("Failed to sign the transaction:",err.Error())
		return
	}
	//fmt.Println(signedTxByte)
	outFile := fmt.Sprintf("signed-%s-%d",fromAddr, rawTx.Nonce)
	err = WriteToFile(outFile, signedTxByte)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("signed transaction is saved in:", outFile)
}
