package tester

import (
    "math/big"
    "testing"

    "github.com/Ankr-network/dccn-common/wallet"
    . "github.com/smartystreets/goconvey/convey"
)

func TestGenerateKeys(t *testing.T) {
    Convey("generate key", t, func() {
        priv, pub, addr := wallet.GenerateKeys()
        address, err := wallet.GetAddressByPublicKey(pub)
        So(addr, ShouldEqual, address)
        _, err = wallet.Sign("123456", priv)
        So(err, ShouldBeNil)
    })
}

func TestSendCoin(t *testing.T) {
    checkoutTestAccount()
    Convey("Test SendCoin main flow, should return nil", t, func() {
        _, _, receiveAddr := wallet.GenerateKeys()
        _, err := wallet.SendCoins(node1, ipPort, adminPrivKey, adminAddress, receiveAddr, sendVal2)
        So(err, ShouldBeNil)
        balanceT, err := wallet.GetBalance(node1, ipPort, receiveAddr)
        So(err, ShouldBeNil)
        So(balanceT, ShouldEqual, "5000000000000000000")
    })

    Convey("TestSendCoin should fail if send more than total supply", t, func() {
        _, _, receiveAddr := wallet.GenerateKeys()
        totalSupply := new(big.Int)
        ankrBase := big.NewInt(1000000000000000000)
        totalSupply = totalSupply.Mul(ankrBase, big.NewInt(1000000000000000000))
        _, err := wallet.SendCoins(node1, ipPort, adminPrivKey, adminAddress, receiveAddr, totalSupply.Add(totalSupply, big.NewInt(1)).String())
        So(err, ShouldBeError)
        _, err = wallet.GetBalance(node1, ipPort, receiveAddr)
        So(err, ShouldBeError)
    })
}

func TestGetMeter(t *testing.T) {
    Convey("test get meter", t, func() {
        _, err := wallet.GetHistoryMetering(node1, ipPort, "datacenter_name", "test-deploy", false, 0, 0)
        So(err, ShouldBeNil)
    })
}

func TestGetBalance(t *testing.T)  {
    Convey("test getbalance", t, func() {
        bal, err := wallet.GetBalance(node1, ipPort, "64BC85F08C03F42B17EAAF5AFFAF9BFAF96CFCB85CA2F3")
        ShouldNotBeNil(err)
        t.Log("64BC85F08C03F42B17EAAF5AFFAF9BFAF96CFCB85CA2F3", bal)
    })
}
