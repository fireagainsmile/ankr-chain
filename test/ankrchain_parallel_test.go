package tester

import (
    "fmt"
    "github.com/Ankr-network/dccn-common/wallet"
    . "github.com/smartystreets/goconvey/convey"
    "math/big"
    "sync"
    "testing"
)

type account struct {
    priv string
    pub  string
    addr string
}

var (
   ipPort       = "26657"     // ip port of nodes
   node1    = "127.0.0.1" //local node ip
   node2    = "127.0.0.1" //
   adminAddress = "B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67"
   adminPrivKey = "wmyZZoMedWlsPUDVCOy+TiVcrIBPcn3WJN8k5cPQgIvC8cbcR10FtdAdzIlqXQJL9hBw1i0RsVjF6Oep/06Ezg==" //node amin private key
   balBase      = big.NewInt(1000000000000000000)                                                            //the base unit of ankr token
   txFee        = big.NewInt(5000000000000000000)                                                            // fee of each transaction
)

//Test multiple account send to once account at the same time
//generate multiple account, set 20 ankr token each
//each account send 10 ankr token to the receive account, check the receive account balance
func TestMultSendOne(t *testing.T) {
   initAccout := account{priv: "wmyZZoMedWlsPUDVCOy+TiVcrIBPcn3WJN8k5cPQgIvC8cbcR10FtdAdzIlqXQJL9hBw1i0RsVjF6Oep/06Ezg==",
       addr: "B508ED0D54597D516A680E7951F18CAD24C7EC9FCFCD67"}
   accounts := generateAccountsFromOne(initAccout, 4)
   sendAmount := new(big.Int)
   sendAmount.Mul(balBase, big.NewInt(10))
   sendToOneAndCheck(accounts, sendAmount.String(), 3)
}

//Test using one account send coins to multiple account at the same time
//generate multiple accounts, init one of them with 100 ankr token
// send 10 ankr token to others from the init account, check balance of the receive accounts and the init account
func TestOneSendMult(t *testing.T) {
    Convey("test multiple account send to one account at the same time", t, func() {
        var accounts []account
        val := "20000000000000000000" //20 ankr token
        for i := 0; i < 20; i++ {
            priv, pub, addr := wallet.GenerateKeys()
            accounts = append(accounts, account{priv: priv, pub: pub, addr: addr})
            _,err := wallet.SendCoins(node1, "26657",adminPrivKey, adminAddress, addr, val)
            ShouldBeNil(err)
        }
        for _, acc := range accounts {
            balance, _ := wallet.GetBalance(node1, ipPort, acc.addr)
            So(balance, ShouldEqual, "15000000000000000000")
        }
        sendVal := "10000000000000000000" //10 ankr token
        for j := 1; j < 20; j++ {
            _, err := wallet.SendCoins(node1, "26657", accounts[j].priv, accounts[j].addr, accounts[0].addr, sendVal)
            So(err, ShouldBeNil)
        }
        balance, err := wallet.GetBalance(node1, "26657", accounts[0].addr)
        So(err, ShouldBeNil)
        So(balance, ShouldEqual, "110000000000000000000")
    })
}

func TestMultSendOnePara(t *testing.T) {
    var accounts []account
    amount := new(big.Int)
    amount.Mul(balBase, big.NewInt(20))
    for i := 0; i < 20; i++ {
        priv, pub, addr := wallet.GenerateKeys()
        accounts = append(accounts, account{priv: priv, pub: pub, addr: addr})
        _,err := wallet.SendCoins(node1, "26657", adminPrivKey,adminAddress,addr, amount.String())
        if err != nil {
            t.Errorf("send coin error %s", err)
        }
        balance, err := wallet.GetBalance(node1, "26657", accounts[0].addr)
        if err != nil {
            t.Errorf("get balance error %s", err)
        }
        fmt.Println("get balance ",balance)
    }
    wg := sync.WaitGroup{}
    txVal := new(big.Int)
    txVal.Mul(balBase, big.NewInt(10))
    for j := 1; j < 20; j++ {
        wg.Add(1)
        go func(j int,wgg *sync.WaitGroup) {
            _, err := wallet.SendCoins(node1, "26657", accounts[j].priv, accounts[j].addr, accounts[0].addr, txVal.String())
            if err != nil {
                t.Errorf("send coins error %s", err)
            }
            wgg.Done()
        }(j,&wg)
    }
    wg.Wait()
    balance, err := wallet.GetBalance(node1, "26657", accounts[0].addr)
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
    var accounts []account
    for i := 0; i < 5; i++ {
        priv, pub, addr := wallet.GenerateKeys()
        amount := new(big.Int)
        amount.Mul(balBase, big.NewInt(20))
        _,err := wallet.SendCoins(node1, "26657", adminPrivKey, adminAddress,addr,amount.String())
        if err != nil {
            t.Errorf("set balance error %s", err)
        }
        accounts = append(accounts, account{priv: priv, pub: pub, addr: addr})
    }
    amount := new(big.Int)
    amount.Mul(balBase, big.NewInt(6))
    for j := 1; j < 5; j++ {
        go func(j int) {
            _, err := wallet.SendCoins(node1, "26657", accounts[j].priv, accounts[j].addr, accounts[0].addr,amount.String() )
            if err != nil {
                fmt.Errorf("send coins error %s", err)
            }
        }(j)
        go func(j int) {
            _, err := wallet.SendCoins(node1, "26657", accounts[j].priv, accounts[j].addr, accounts[0].addr, amount.String())
            if err != nil {
                fmt.Errorf("send coins error %s", err)
            }
        }(j)
    }
    balance, err := wallet.GetBalance(node1, "26657", accounts[0].addr)
    if err != nil {
        t.Errorf("get balance error %s", err)
    }
    fmt.Printf("%s", balance)
}

func TestOneDiffNodeSendMult(t *testing.T) {
    var accounts []account
    for i := 0; i < 2; i++ {
        priv, pub, addr := wallet.GenerateKeys()
        accounts = append(accounts, account{priv: priv, pub: pub, addr: addr})
    }
    amount := new(big.Int)
    amount.Mul(balBase, big.NewInt(50))
    _,err := wallet.SendCoins(node1, "26657", adminPrivKey,adminAddress,accounts[0].addr, amount.String())
    if err != nil {
        t.Errorf("setBalance error %s", err.Error())
    }
    var j int
    txVal := new(big.Int)
    txVal.Mul(balBase, big.NewInt(6))
    for j = 1; j < 2; j++ {
        go func(j int) {
            _, err := wallet.SendCoins(node1, "26657", accounts[0].priv, accounts[0].addr, accounts[j].addr, txVal.String())
            if err != nil {
                fmt.Errorf("snedCoins error %s", err)
            }
        }(j)
        go func(j int) {
            _, err := wallet.SendCoins(node1, "26657", accounts[0].priv, accounts[0].addr, accounts[j].addr, txVal.String())
            if err != nil {
                fmt.Errorf("snedCoins error %s", err)
            }
        }(j)
    }
    balance, err := wallet.GetBalance(node1, "26657", accounts[0].addr)
    fmt.Printf("%s", balance)
}

func sendToOneAndCheck(accounts []account, value string, numTxs int) {
    _, _, addr := wallet.GenerateKeys()
    var wg sync.WaitGroup
    for _, acc := range accounts {
        wg.Add(1)
        go func(a account) {
            for i := 0; i < numTxs; i++ {
                _, err := wallet.SendCoins(node1, ipPort, a.priv, a.addr, addr, value)
                if err != nil {
                    panic(err)
                }
            }
            wg.Done()
        }(acc)
    }
    //check balance of the receiver
    totalReceive := new(big.Int)
    txValue, _ := new(big.Int).SetString(value, 10)
    txValue.Sub(txValue, txFee)
    totalReceive.Mul(txValue, big.NewInt(int64(numTxs)))
    totalReceive.Mul(totalReceive, big.NewInt(int64(len(accounts))))
    wg.Wait()
    bal, err := wallet.GetBalance(node1, ipPort, addr)
    if err != nil {
        panic(err)
    }
    if bal != totalReceive.String() {
        panic(fmt.Errorf(fmt.Sprintf("Expect %s, got %s", totalReceive.String(), bal)))
    }
}

func generateAccountsFromOne(initAccount account, num uint) []account {
    accounts := make([]account, 0, 1 << num)
    accounts = append(accounts,initAccount)
    numRoutine := 4
    oldAccountChan := make(chan account, numRoutine)
    newAccountChan := make(chan account, numRoutine)
    //stopChannels := make([]chan struct{}, 4)
    for i := 0; i< numRoutine; i++ {
        go generateAcconuntLoop(oldAccountChan, newAccountChan)
    }

    var wg sync.WaitGroup
    for i := 0; i< int(num); i++ {
        wg.Add(1)
        go func(length int) {
            for {
                select {
                case newAcc := <- newAccountChan:
                    accounts = append(accounts, newAcc)
                    length --
                    if length <= 0{
                        wg.Done()
                        return
                    }
                }
            }
        }(len(accounts))

        for _, acc := range accounts {
            oldAccountChan <- acc
        }
        wg.Wait()
    }
    return accounts
}


func generateAcconuntLoop(oldAccount <- chan account, newAccount chan<- account){
    for  {
        select {
        case acc := <-oldAccount:
            priv, pub, addr := wallet.GenerateKeys()
            bal, _ := wallet.GetBalance(node1, ipPort, acc.addr)
            banlance, _ := new(big.Int).SetString(bal, 10)
            banlance.Div(banlance, big.NewInt(2))
            wallet.SendCoins(node1, ipPort, acc.priv, acc.addr, addr, banlance.String())
            newAccount <- account{priv, pub, addr }
        }
    }
}