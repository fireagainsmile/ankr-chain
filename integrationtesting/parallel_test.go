package tester

import (
    "fmt"
    "math/big"
    "os"
    "sync"
    "testing"

    "github.com/Ankr-network/ankr-chain/consensus"
    "github.com/Ankr-network/dccn-common/wallet"
    . "github.com/smartystreets/goconvey/convey"
)

type account struct {
    priv string
    pub  string
    addr string
}

var (
    ipPort       = "26657"                     // ip port of nodes
    //node1        = os.Getenv("ANKR_URL")
    node1 = "chain-dev.dccn.ankr.com"
    node2        = os.Getenv("ANKR_URL")
    adminAddress = "B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67"
    adminPrivKey = "wmyZZoMedWlsPUDVCOy+TiVcrIBPcn3WJN8k5cPQgIvC8cbcR10FtdAdzIlqXQJL9hBw1i0RsVjF6Oep/06Ezg==" //node amin private key
    sendVal      = "20000000000000000000" //20 ankr token
    sendVal2     = "10000000000000000000" // 10 ankr token
    expectVal2   = "15000000000000000000"
    leftVal = "5000000000000000000"
    balBase = "1000000000000000000"
)

func checkoutTestAccount() {
    fmt.Println("test node ",node1)
    if ankrchain.ADMIN_OP_FUND_PUBKEY == "cGSgVIfAsXWbuWImGxJlNzfqruzuGA+4JXv5gfB0FyY=" {
        return
    }
}

// send coins, check single transaction logic
func TestSendSingleTransaction(t *testing.T)  {
    base, _ := new(big.Int).SetString(balBase, 10)
    initPriv, _, initAddr := wallet.GenerateKeys()
    initAmount := new(big.Int)
    initAmount = initAmount.Mul(base, big.NewInt(20))
    _, err := wallet.SendCoins(node1, ipPort, adminPrivKey, adminAddress, initAddr, initAmount.String())
    if err != nil {
        t.Error(err)
    }
    receivePriv, _, receiveAddr := wallet.GenerateKeys()
    sendAmount := new(big.Int)
    leftAmount := initAmount.Mul(base, big.NewInt(15))

    Convey("test send more than account balance", t, func() {
      sendAmount = sendAmount.Mul(base, big.NewInt(30))
      _, err = wallet.SendCoins(node1, ipPort, initPriv, initAddr, receiveAddr, sendAmount.String())
      So(err, ShouldBeError)
      bal, err := wallet.GetBalance(node1, ipPort, initAddr)
      So(err, ShouldBeNil)
      So(bal, ShouldEqual, leftAmount.String())
      _, err = wallet.GetBalance(node1, ipPort, receiveAddr)
      So(err, ShouldBeError)
    })

    Convey("test send all the amount of the account", t, func() {
        sendAmount = sendAmount.Mul(base, big.NewInt(15))
        _, err = wallet.SendCoins(node1, ipPort, initPriv, initAddr, receiveAddr, sendAmount.String())
        So(err, ShouldBeNil)
        bal, err := wallet.GetBalance(node1, ipPort, initAddr)
        So(err, ShouldBeNil)
        So(bal, ShouldEqual, "0")
        bal, err = wallet.GetBalance(node1, ipPort, receiveAddr)
        So(err, ShouldBeNil)
        leftAmount = leftAmount.Mul(base, big.NewInt(10))
        So(bal, ShouldEqual, leftAmount.String())
    })

    Convey("test send back", t, func() {
        sendAmount = sendAmount.Mul(base, big.NewInt(10))
        _, err = wallet.SendCoins(node1, ipPort, receivePriv, receiveAddr, initAddr, sendAmount.String())
        So(err, ShouldBeNil)
        bal, err := wallet.GetBalance(node1, ipPort, receiveAddr)
        So(err, ShouldBeNil)
        So(bal, ShouldEqual, "0")
        bal, err = wallet.GetBalance(node1, ipPort, initAddr)
        So(err, ShouldBeNil)
        leftAmount = leftAmount.Mul(base, big.NewInt(5))
        So(bal, ShouldEqual, leftAmount.String())
    })
}

// test multiple account send to one account at the same time
func TestMultSendOne(t *testing.T) {
    checkoutTestAccount()
    var accounts []account
    gasAccount := "64BC85F08C03F42B17EAAF5AFFAF9BFAF96CFCB85CA2F3"
    gasBal, _ := wallet.GetBalance(node1, ipPort, gasAccount)
    t.Log(gasAccount,gasBal)
    for i := 0; i < 6; i++ {
        priv, pub, addr := wallet.GenerateKeys()
        _, err := wallet.SendCoins(node1, ipPort, adminPrivKey, adminAddress, addr, sendVal)
        ShouldBeNil(err)
        accounts = append(accounts, account{priv: priv, pub: pub, addr: addr})
        bal, err := wallet.GetBalance(node1, ipPort, addr)
        ShouldBeNil(err)
        ShouldEqual(bal, expectVal2)
    }
    var wg sync.WaitGroup
    _, _, receiveAddr := wallet.GenerateKeys()
    for j := 0; j < 6; j++ {
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
            ShouldEqual(bal, leftVal)
            wg.Done()
        }(j)
    }
    wg.Wait()
    gasBal, _ = wallet.GetBalance(node1, ipPort, gasAccount)
    t.Log(gasAccount,gasBal)
    bal, err := wallet.GetBalance(node1, ipPort, receiveAddr)
    if err != nil {
        t.Errorf("TestMultSendOne GetBalance error %s", err)
    }
    if bal != "30000000000000000000" {
        t.Errorf("test MultSendOne get balance error %s", bal)
    }
}

// test using one account continuously send to other account
func TestOneKeepSending(t *testing.T) {
    checkoutTestAccount()
    Convey("test using one account continuously send to different account", t, func() {
        base, _ := new(big.Int).SetString(balBase, 10)

        initAmount := new(big.Int)
        initAmount = initAmount.Mul(base, big.NewInt(110))
        //init account
        initPriv, _, initAddr := wallet.GenerateKeys()
         _, err := wallet.SendCoins(node1, ipPort, adminPrivKey, adminAddress, initAddr, initAmount.String())
         So(err, ShouldBeNil)

        sendVal := new(big.Int)
        sendVal = sendVal.Mul(base, big.NewInt(10))
        var accounts []account
        for i := 0; i < 10; i++ {
            priv, pub, addr := wallet.GenerateKeys()
            accounts = append(accounts, account{priv: priv, pub: pub, addr: addr})
            _, err := wallet.SendCoins(node1, ipPort, initPriv, initAddr, addr, sendVal.String())
            ShouldBeNil(err)
        }
        leftBal := initAmount.Mul(base, big.NewInt(5))
        initBal, _ := wallet.GetBalance(node1, ipPort, initAddr)
        So(initBal, ShouldEqual, leftBal.String())

        receiveBal := initAmount.Mul(base, big.NewInt(5))
        for _, acc := range accounts {
            bal, _ := wallet.GetBalance(node1, ipPort, acc.addr)
            So(bal, ShouldEqual, receiveBal.String())
        }
    })

    Convey("test using one account continuously sending to another account", t, func() {
        base, _ := new(big.Int).SetString(balBase, 10)
        initAmount := new(big.Int)
        initAmount = initAmount.Mul(base, big.NewInt(110))
        //init account
        initPriv, _, initAddr := wallet.GenerateKeys()
        _, err := wallet.SendCoins(node1, ipPort, adminPrivKey, adminAddress, initAddr, initAmount.String())
        So(err, ShouldBeNil)
        sendVal := new(big.Int)
        sendVal = sendVal.Mul(base, big.NewInt(10))
        _, _, receiveAddr := wallet.GenerateKeys()
        for i := 0; i < 10 ; i ++ {
            _, err := wallet.SendCoins(node1, ipPort, initPriv, initAddr, receiveAddr, sendVal.String())
            So(err, ShouldBeNil)
        }
        checkVal := sendVal.Mul(base, big.NewInt(5))
        initAccBal, _ := wallet.GetBalance(node1, ipPort, initAddr)
        So(initAccBal, ShouldEqual, checkVal.String())

        receiveVal := sendVal.Mul(base, big.NewInt(50))
        receiveBal, _ := wallet.GetBalance(node1, ipPort, receiveAddr)
        So(receiveBal, ShouldEqual, receiveVal.String())
    })
}

// todo
// need another url
func TestMultNodeSendOnePara(t *testing.T) {
    checkoutTestAccount()
    var accounts []account
    for i := 0; i < 5; i++ {
        priv, pub, addr := wallet.GenerateKeys()
        _, err := wallet.SendCoins(node1, ipPort, adminPrivKey, adminAddress, addr, sendVal)
        if err != nil {
            t.Errorf("set balance error %s", err)
        }
        accounts = append(accounts, account{priv: priv, pub: pub, addr: addr})
    }
    var sendCoin = "6000000000000000000"
    node := node1
    wg := sync.WaitGroup{}
    for j := 1; j < 5; j++ {
        wg.Add(1)
        go func(j int, node string, wgg *sync.WaitGroup) {
            if j%2 != 0 {
                node = node2
            }
            _, err := wallet.SendCoins(node, ipPort, accounts[j].priv, accounts[j].addr, accounts[0].addr, sendCoin)
            if err != nil {
                fmt.Errorf("send coins error %s", err)
            }
            wgg.Done()
        }(j, node, &wg)
    }
    wg.Wait()
    balance, err := wallet.GetBalance(node1, ipPort, accounts[0].addr)
    if err != nil {
        t.Errorf("get balance error %s", err)
    }
    if balance != "19000000000000000000" {
        t.Errorf("get balance is not %s", balance)
    }
}

//todo
// need another node url
func TestOneDiffNodeSendMult(t *testing.T) {
    checkoutTestAccount()
    var accounts []account
    for i := 0; i < 5; i++ {
        priv, pub, addr := wallet.GenerateKeys()
        accounts = append(accounts, account{priv: priv, pub: pub, addr: addr})
    }
    _, err := wallet.SendCoins(node1, ipPort, adminPrivKey, adminAddress, accounts[0].addr, sendVal)
    if err != nil {
        t.Errorf("setBalance error %s", err.Error())
    }
    var j int
    var sendCoin = "6000000000000000000"
    var node = node1
    var wg = sync.WaitGroup{}
    for j = 1; j < 4; j++ {
        wg.Add(1)
        if j%2 == 0 {
            node = node2
        }
        go func(j int, node string, wgg *sync.WaitGroup) {
            _, err = wallet.SendCoins(node, ipPort, accounts[0].priv, accounts[0].addr, accounts[j].addr, sendCoin)
            wgg.Done()
        }(j, node, &wg)
    }
    wg.Wait()
    balance, err := wallet.GetBalance(node, ipPort, accounts[0].addr)
    if balance != "9000000000000000000" {
        t.Errorf("OneDiffNodeSendMult get balance is not equal %s", balance)
    }

}
