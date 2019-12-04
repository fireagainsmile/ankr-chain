package cmd

import (
	"fmt"
	"github.com/agiledragon/gomonkey"
	"github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"testing"
)

//normal cases
func TestGetBalance(t *testing.T) {
	convey.Convey("test get balance", t, func() {

		args := []string{"query", "balance", "--address", "95CD00025C3807CEE9804D19B1E410A30A47B303371C12", "--nodeurl", localUrl}
		cmd := RootCmd
		cmd.SetArgs(args)
		err := cmd.Execute()
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestGenAccount(t *testing.T) {
	convey.Convey("test generate account command", t, func() {
		args := []string{"account", "genaccount"}
		cmd := RootCmd
		cmd.SetArgs(args)
		err := cmd.Execute()
		fmt.Println(err)
		//convey.So(err, convey.ShouldBeNil)
		//removeKeyStore("./tmp", address)
	})
}

var keyStoreString = `{"address":"B47982CF51CD7718FE25BA96B707C449F4F917949E7A25","publickey":"TokCLsXk6Z4lhyr/v7PkSksuHkDQURrjVn1Pt8JsKNI=","crypto":{"cipher":"aes-128-ctr","ciphertext":"3e0d4968e664d201682fecbba66e9fa2b8f5d4ccd9445c45705a9864743234aed4e2489d26025fd97a3fa7f23b3714e580b7a38bac1b5c734e75a6c5fac9cd922239f54934f7cc073b65f5fb914166f7964f8fc07093a4c6","cipherparams":{"iv":"e9b282f254f2bcdf98e42cb32c127b06"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"f3698a69cb65c64fa9548febbd19174f1af898cd56590420e5a5fe0578bd351c"},"mac":"eaf9d07ff51732eb3c58e561f0ffc37bf099c5570d2a51388cc5b2f3d7420860"},"version":3}`
func TestExportPriv(t *testing.T) {
	convey.Convey("test exporting private key from keystore", t, func() {
		args := []string{"account", "exportprivatekey", "--file", "./keystore"}
		cmd := RootCmd
		cmd.SetArgs(args)
		err := cmd.Execute()
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestResetPWD(t *testing.T) {
	convey.Convey("test exporting private key from keystore", t, func() {
		filePatch := gomonkey.ApplyFunc(ioutil.ReadFile, func(file string) ([]byte, error){
			return []byte(keyStoreString), nil
		})
		defer filePatch.Reset()

		args := []string{"account", "resetpwd", "--file", "tmp/keystore"}
		cmd := RootCmd
		cmd.SetArgs(args)
		err := cmd.Execute()
		convey.So(err, convey.ShouldBeNil)
	})
}

