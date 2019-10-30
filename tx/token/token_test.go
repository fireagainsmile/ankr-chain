package token

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/Ankr-network/dccn-common/wallet"
	"github.com/stretchr/testify/assert"
)

var (
	ipPort       = "26657"                     // ip port of nodes
	//node1        = os.Getenv("ANKR_URL")
	node1 = "127.0.0.1"
	node2        = os.Getenv("ANKR_URL")
	adminAddress = "B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67"
	adminPrivKey = "wmyZZoMedWlsPUDVCOy+TiVcrIBPcn3WJN8k5cPQgIvC8cbcR10FtdAdzIlqXQJL9hBw1i0RsVjF6Oep/06Ezg==" //node amin private key
	sendVal      = "20000000000000000000" //20 ankr token
	sendVal2     = "10000000000000000000" // 10 ankr token
	expectVal2   = "15000000000000000000"
	leftVal = "5000000000000000000"
)

type accountT struct {
	priv string
	pub string
	addr string
}

func TestMultSendOne(t *testing.T) {

	var accounts []accountT
	for i := 0; i < 2; i++ {
		priv, pub, addr := wallet.GenerateKeys()
		_, err := wallet.SendCoins(node1, ipPort, adminPrivKey, adminAddress, addr, sendVal)
		assert.Equal(t, err, nil)
		accounts = append(accounts, accountT{priv: priv, pub: pub, addr: addr})
		bal, err := wallet.GetBalance(node1, ipPort, addr)
		assert.Equal(t, err, nil)
		//assert.Equal(t, bal, expectVal2)
		fmt.Printf("%s bal:%s\n", addr, bal)
	}
	var wg sync.WaitGroup
	_, _, receiveAddr := wallet.GenerateKeys()
	for j := 0; j < 2; j++ {
		wg.Add(1)
		go func(j int) {
			_, err := wallet.SendCoins(node1, ipPort, accounts[j].priv, accounts[j].addr, receiveAddr, sendVal2)
			if err != nil {
				t.Errorf("TestMultSendOne send error %s", err)
			}
			bal, err := wallet.GetBalance(node1, ipPort, accounts[j].addr)
			if err != nil {
				t.Errorf("Getbalance error %s", err)
			}
			//assert.Equal(t, bal, leftVal)

			fmt.Printf("%s bal:%s\n", accounts[j].addr, bal)

			fmt.Printf("receive %s bal:%s\n", receiveAddr, bal)

			wg.Done()
		}(j)
	}
	wg.Wait()
	bal, err := wallet.GetBalance(node1, ipPort, receiveAddr)
	if err != nil {
		t.Errorf("TestMultSendOne GetBalance error %s", err)
	}
	//if bal != "30000000000000000000" {
	//   t.Errorf("test MultSendOne get balance error %s", bal)
	// }
	fmt.Printf("final receive %s bal:%s\n", receiveAddr, bal)
}

