package tester

import (
    "context"
    "fmt"
    "testing"
    "time"

    "github.com/Ankr-network/dccn-common/wallet"
    "github.com/tendermint/tendermint/libs/pubsub/query"
    "github.com/tendermint/tendermint/types"
)

var waitForEventTimeout = 30 * time.Second

func TestTxAndAddress(t *testing.T) {
   Hclient := NewHttpClient()
   _, _, addr := wallet.GenerateKeys()
   qStr := fmt.Sprintf("tm.event='%s' AND app.toaddress CONTAINS '%s'", types.EventTx, addr)
   q := query.MustParse(qStr)
   ctx, cancel := context.WithTimeout(context.Background(), waitForEventTimeout)
   defer cancel()
   eventChTx, err := Hclient.Subscribe(ctx, "helper", q.String())
   if err != nil {
       t.Error("sub tx and toaddress error", err)
   }
   go func() {
       for {
           select {
           case evt := <-eventChTx:
               evDataTx := evt.Data.(types.EventDataTx)
               if evDataTx.Tx != nil {
                   t.Log(evDataTx.Height)
               } else {
                   t.Error("TestTxAndAddress error")
               }
               err:=Hclient.Unsubscribe(ctx,"helper", q.String())
               if err !=nil{
                   t.Error("TestTxAndAddress Unsubscribe error")
               }
               return
           }
       }
   }()
   wallet.SendCoins(node1, ipPort, adminPrivKey, adminAddress, addr, "1000000000000000000")
   time.Sleep(15 *time.Second)
}

func TestTxAndsubMuitAddress(t *testing.T) {
    Hclient := NewHttpClient()
    _, _, addr1 := wallet.GenerateKeys()
    _, _, addr2 := wallet.GenerateKeys()
    addr:=fmt.Sprintf("%s,%s",addr2,addr1)
    qStr := fmt.Sprintf("tm.event='%s' AND app.toaddress CONTAINS '%s'", types.EventTx, addr)
    q := query.MustParse(qStr)
    ctx, cancel := context.WithTimeout(context.Background(), waitForEventTimeout)
    defer cancel()
    eventChTx, err := Hclient.Subscribe(ctx, "helper", q.String())
    if err != nil {
        t.Error("sub tx and toaddress error", err)
    }
    go func() {
        for {
            select {
            case evt := <-eventChTx:
                evDataTx := evt.Data.(types.EventDataTx)
                if evDataTx.Tx != nil {
                    t.Log(evDataTx.Height)
                } else {
                    t.Error("TestTxAndAddress error")
                }
                err:=Hclient.Unsubscribe(ctx,"helper", q.String())
                if err!=nil{
                    t.Error("TestTxAndsubMuitAddress Unsubscribe error")
                }
                return
            }
        }
    }()
    wallet.SendCoins(node1, ipPort, adminPrivKey, adminAddress, addr, "1000000000000000000")
    time.Sleep(15 *time.Second)
}

func TestSubEventTx(t *testing.T) {
   Hclient := NewHttpClient()
   ctx, cancel := context.WithTimeout(context.Background(), waitForEventTimeout)
   defer cancel()
   eventCh, err := Hclient.Subscribe(ctx, "helper", "tm.event='Tx'")
   if err != nil {
       t.Error("test sub tx error", err)
   }
   go func() {
       for {
           select {
           case evt := <-eventCh:
               evDataTx := evt.Data.(types.EventDataTx)
               t.Log(evDataTx.Height)
               err:=Hclient.Unsubscribe(ctx,"helper", "tm.event='Tx'")
               if err!=nil{
                   t.Error("TestSubEventTx Unsubscribe error")
               }
               //So(evDataTx.Result.Code, ShouldEqual, 0)
               return
           }
       }
   }()
   time.Sleep(15 *time.Second)
}

func TestSubNewBlock(t *testing.T) {
   Hclient := NewHttpClient()
   ctx, cancel := context.WithTimeout(context.Background(), waitForEventTimeout)
   defer cancel()
   eventCh, err := Hclient.Subscribe(ctx, "helper", "tm.event='NewBlock'")
   if err != nil {
       t.Error("test sub new block error", err)
   }
   go func() {
       _, err := Hclient.Status()
       if err == nil {
           for {
               select {
               case evt := <-eventCh:
                   evData := evt.Data.(types.EventDataNewBlock)
                   if evData.Block.Height != 0 {
                       fmt.Printf("%d", evData.Block.Height)
                   } else {
                       t.Error("NewBlock error")
                   }
                   return
               }
           }
       } else {
           t.Error("stat SyncInfo ", err)
           return
       }
   }()
   time.Sleep(15 *time.Second)

}

func TestSubNewBlockHeader(t *testing.T) {
   Hclient := NewHttpClient()
   ctx, cancel := context.WithTimeout(context.Background(), waitForEventTimeout)
   defer cancel()
   eventCh, err := Hclient.Subscribe(ctx, "helper", "tm.event='NewBlockHeader'")
   if err != nil {
       t.Error("test sub new block error", err)
   }
   go func() {
       stat, err := Hclient.Status()
       if err == nil {
           for {
               select {
               case evt := <-eventCh:
                   evData := evt.Data.(types.EventDataNewBlockHeader)
                   if evData.Header.Height != 0 {
                       fmt.Printf("%d", evData.Header.Height)
                   } else {
                       t.Error("NewBlock error")
                   }
                   return
               }
           }
       } else {
           fmt.Printf("%d", stat.SyncInfo.LatestBlockHeight)
           //So(stat.SyncInfo.LatestBlockHeight, ShouldBeError)
           return
       }
   }()
   time.Sleep(15 *time.Second)
}

func TestSubEventVote(t *testing.T) {
   Hclient := NewHttpClient()
   ctx, cancel := context.WithTimeout(context.Background(), waitForEventTimeout)
   defer cancel()
   eventCh, err := Hclient.Subscribe(ctx, "helper", "tm.event='vote'")
   if err != nil {
       t.Error("event vote", err)
   }
   go func() {
       for {
           select {
           case evt := <-eventCh:
               evData := evt.Data.(types.EventDataVote)
               fmt.Printf("%d", evData.Vote.Height)
               //So(evData.Vote.Height, ShouldNotBeEmpty)
               return
           }
       }
   }()
   time.Sleep(15 *time.Second)
}