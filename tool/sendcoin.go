package main

import (
	"fmt"
	"github.com/Ankr-network/dccn-common/wallet"
)

/*
test data(priv key, pub key, address):
wmyZZoMedWlsPUDVCOy+TiVcrIBPcn3WJN8k5cPQgIvC8cbcR10FtdAdzIlqXQJL9hBw1i0RsVjF6Oep/06Ezg==
wvHG3EddBbXQHcyJal0CS/YQcNYtEbFYxejnqf9OhM4=
B508ED0D54597D516A680E7951F18CAD24C7EC9F
*/

/*
   test API to send 88 tokens from address 1 to address 2
   address 1:B508ED0D54597D516A680E7951F18CAD24C7EC9F
   address 2:0D9FE6A785C830D2BE66FE40E0E7FE3D9838456C
*/

//NjkJiEleYz1nShwlYbQxEhnb3SyPSlgyLY6t32lubM/p37BkgDL0UAS0dM1z9QnJFtRweK9lXKb29ZrLWAjekQ==
//6d+wZIAy9FAEtHTNc/UJyRbUcHivZVym9vWay1gI3pE=
//5AEBA6EB8BC51DA277CCF1EF229F0C05D9535FA36CC872
func main() {
	hash, err := wallet.SendCoins("chain-dev.dccn.ankr.com", "26657",                                             //"chain-dev.dccn.ankr.com"
		"wmyZZoMedWlsPUDVCOy+TiVcrIBPcn3WJN8k5cPQgIvC8cbcR10FtdAdzIlqXQJL9hBw1i0RsVjF6Oep/06Ezg==", //priv_key
		"B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67",                                           //from address
		"5AEBA6EB8BC51DA277CCF1EF229F0C05D9535FA36CC872",                                           //to address
		"6000000000000000000")                                                                    //amount
	//"wvHG3EddBbXQHcyJal0CS/YQcNYtEbFYxejnqf9OhM4=") //pub key not needed

	if err == nil {
		fmt.Println("send success.", hash)
	} else {
		fmt.Println("send failure, ", err)
	}
}