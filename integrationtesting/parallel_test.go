package tester

import (
    "fmt"
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
    node1        = "127.0.0.1" //"chain-stage.dccn.ankr.com" //local node ip
    node2        = "127.0.0.1" //"chain-stage.dccn.ankr.com" //
    adminAddress = "B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67"
    adminPrivKey = "wmyZZoMedWlsPUDVCOy+TiVcrIBPcn3WJN8k5cPQgIvC8cbcR10FtdAdzIlqXQJL9hBw1i0RsVjF6Oep/06Ezg==" //node amin private key
    sendVal      = "20000000000000000000"
    sendVal2     = "10000000000000000000"
    expectVal2   = "15000000000000000000"
)

func checkoutTestAccount() {
    if ankrchain.ADMIN_OP_FUND_PUBKEY == "cGSgVIfAsXWbuWImGxJlNzfqruzuGA+4JXv5gfB0FyY=" {
        return
    }
}
func TestMultSendOne(t *testing.T) {
    checkoutTestAccount()
    var accounts []account
    for i := 0; i < 6; i++ {
        priv, pub, addr := wallet.GenerateKeys()
        _, err := wallet.SendCoins(node1, ipPort, adminPrivKey, adminAddress, addr, sendVal)
        ShouldBeNil(err)
        accounts = append(accounts, account{priv: priv, pub: pub, addr: addr})
        bal, err := wallet.GetBalance(node1, ipPort, accounts[0].addr)
        ShouldBeNil(err)
        ShouldNotEqual(bal, expectVal2)
    }
    var wg sync.WaitGroup
    for j := 1; j < 6; j++ {
        wg.Add(1)
        go func(j int, wgg *sync.WaitGroup) {
            _, err := wallet.SendCoins(node1, ipPort, accounts[j].priv, accounts[j].addr, accounts[0].addr, sendVal2)
            if err != nil {
                t.Errorf("TestMultSendOne send error %s", err)
            }
            wgg.Done()
        }(j, &wg)
    }
    wg.Wait()
    bal, err := wallet.GetBalance(node1, ipPort, accounts[0].addr)
    if err != nil {
        t.Errorf("TestMultSendOne GetBalance error %s", err)
    }
    if bal != "40000000000000000000" {
        t.Errorf("test MultSendOne get balance error %s", err)
    }
}

//Test using one account send coins to multiple account at the same time
//generate multiple accounts, init one of them with 100 ankr token
// send 10 ankr token to others from the init account, check balance of the receive accounts and the init account
func TestOneSendMult(t *testing.T) {
    checkoutTestAccount()
    Convey("test multiple account send to one account at the same time", t, func() {
        var accounts []account
        val := "20000000000000000000" //20 ankr token
        for i := 0; i < 20; i++ {
            priv, pub, addr := wallet.GenerateKeys()
            accounts = append(accounts, account{priv: priv, pub: pub, addr: addr})
            _, err := wallet.SendCoins(node1, ipPort, adminPrivKey, adminAddress, addr, val)
            ShouldBeNil(err)
        }
        for _, acc := range accounts {
            balance, _ := wallet.GetBalance(node1, ipPort, acc.addr)
            So(balance, ShouldEqual, "15000000000000000000")
        }
        sendVal := "10000000000000000000" //10 ankr token
        for j := 1; j < 20; j++ {
            _, err := wallet.SendCoins(node1, ipPort, accounts[j].priv, accounts[j].addr, accounts[0].addr, sendVal)
            So(err, ShouldBeNil)
        }
        balance, err := wallet.GetBalance(node1, ipPort, accounts[0].addr)
        So(err, ShouldBeNil)
        So(balance, ShouldEqual, "110000000000000000000")
    })
}

func TestMultSendOnePara(t *testing.T) {
    checkoutTestAccount()
    var accounts []account
    for i := 0; i < 20; i++ {
        priv, pub, addr := wallet.GenerateKeys()
        accounts = append(accounts, account{priv: priv, pub: pub, addr: addr})
        _, err := wallet.SendCoins(node1, ipPort, adminPrivKey, adminAddress, addr, sendVal)
        if err != nil {
            t.Errorf("send coin error %s", err)
        }
        balance, err := wallet.GetBalance(node1, ipPort, accounts[0].addr)
        if err != nil {
            t.Errorf("get balance error %s", err)
        }
        fmt.Println("get balance ", balance)
    }
    wg := sync.WaitGroup{}
    for j := 1; j < 20; j++ {
        wg.Add(1)
        go func(j int, wgg *sync.WaitGroup) {
            _, err := wallet.SendCoins(node1, ipPort, accounts[j].priv, accounts[j].addr, accounts[0].addr, sendVal2)
            if err != nil {
                t.Errorf("send coins error %s", err)
            }
            wgg.Done()
        }(j, &wg)
    }
    wg.Wait()
    balance, err := wallet.GetBalance(node1, ipPort, accounts[0].addr)
    if err != nil {
        t.Errorf("get balance error %s", err)
    }
    if balance == "110000000000000000000" {
        fmt.Printf("%s", balance)
    } else {
        t.Errorf("get balance amount error %s ", err)
    }
}

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
