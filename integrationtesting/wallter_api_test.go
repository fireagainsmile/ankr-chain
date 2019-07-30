package tester

import (
    "github.com/Ankr-network/dccn-common/wallet"
    . "github.com/smartystreets/goconvey/convey"
    "math"
    "math/big"
    "testing"
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

func TestSetBalance(t *testing.T) {
    sendAmount       := "100000000000000000000"
  Convey("Test SetBalance main flow, should return nil", t, func() {
      _, _, addr := wallet.GenerateKeys()
      wallet.SendCoins(node1, ipPort,adminPrivKey,adminAddress,addr,sendAmount)
      err := wallet.SetBalance(node1, ipPort, adminAddress, sendAmount, adminPrivKey)
      ShouldBeNil(err)

      //check SetBalance result
      balance, err := wallet.GetBalance(node1, ipPort, addr)
      So(err, ShouldBeNil)
      So(balance, ShouldEqual, "95000000000000000000")
  })

  //test SetBalance when the value is set a neg number
  Convey("TestSetBalance should fail if set balance a nge number", t, func() {
      setVal := big.NewInt(-1)
      _, _, addr := wallet.GenerateKeys()
      err := wallet.SetBalance(node1, ipPort, addr, setVal.String(), adminPrivKey)
      So(err, ShouldBeError)
  })

}
//
func TestSendCoin(t *testing.T) {
  //send coin main flow, normal case
  //init send address with 10 ankr token, send all to a receive address
  //check balance of both address
    _, _, receiveAddr := wallet.GenerateKeys()
  Convey("Test SendCoin main flow, should return nil", t, func() {
      _, err := wallet.SendCoins(node1, ipPort, adminPrivKey, adminAddress, receiveAddr, sendVal2)
      So(err, ShouldBeNil)
      balanceT, err := wallet.GetBalance(node1, ipPort, receiveAddr)
      So(err, ShouldBeNil)
      So(balanceT, ShouldEqual, "5000000000000000000")
  })

  Convey("TestSendCoin should success if send more than max uint64", t, func() {
      maxmum := new(big.Int).SetUint64(math.MaxUint64)
      _, err := wallet.SendCoins(node1, ipPort, adminPrivKey, adminAddress, receiveAddr, maxmum.Add(maxmum, big.NewInt(1)).String())
      So(err, ShouldBeNil)
      balance1, err := wallet.GetBalance(node1, ipPort, receiveAddr)
      So(err, ShouldBeNil)
      So(balance1, ShouldEqual, "18446744073709551616")
  })

  Convey("TestSendCoin should fail if send more than total supply", t, func() {
      totalSupply := new(big.Int)
      ankrBase := big.NewInt(1000000000000000000)
      totalSupply = totalSupply.Mul(ankrBase, big.NewInt(1000000000000000000))
      _, err := wallet.SendCoins(node1, ipPort, adminPrivKey, adminAddress, receiveAddr, totalSupply.Add(totalSupply, big.NewInt(1)).String())
      So(err, ShouldBeError)
      balance1, err := wallet.GetBalance(node1, ipPort, receiveAddr)
      So(err, ShouldBeNil)
      So(balance1, ShouldEqual, "18446744073709551616")
  })
}

func TestGetMeter(t *testing.T) {
  Convey("test get meter", t, func() {
      _, err := wallet.GetHistoryMetering(node1, ipPort, "datacenter_name", "test-deploy", false, 0, 0)
      So(err, ShouldBeNil)
  })
}
