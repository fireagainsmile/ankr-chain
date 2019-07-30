package tester

import (
    "testing"
    "context"
    "fmt"
    "time"

   "github.com/Ankr-network/dccn-common/wallet"
   . "github.com/smartystreets/goconvey/convey"
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
               //So(evDataTx.Result.Code, ShouldEqual, 0)
               return
           }
       }
   }()
   wallet.SendCoins(node1, ipPort, adminPrivKey, adminAddress, addr, "1000000000000000000")
   time.Sleep(10 *time.Second)
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
               //So(evDataTx.Result.Code, ShouldEqual, 0)
               return
           }
       }
   }()
   time.Sleep(10 *time.Second)
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
   time.Sleep(10 *time.Second)

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
   time.Sleep(10 *time.Second)
}

func TestSubEventLock(t *testing.T) {
   Hclient := NewHttpClient()
   ctx, cancel := context.WithTimeout(context.Background(), waitForEventTimeout)
   defer cancel()
   eventCh, err := Hclient.Subscribe(ctx, "helper", "tm.event='Lock'")
   if err != nil {
       t.Error("test sub lock error", err)
   }
   go func() {
       stat, err := Hclient.Status()
       if err == nil {
           for {
               select {
               case evt := <-eventCh:
                   evData := evt.Data.(types.EventDataRoundState)
                   fmt.Printf("%d", evData.Height)
                   //So(evData.Height, ShouldBeGreaterThan, stat.SyncInfo.LatestBlockHeight)
                   return
               }
           }
       } else {
           So(stat, ShouldBeError)
           return
       }
   }()
   time.Sleep(10 *time.Second)
}

func TestSubEventUnLock(t *testing.T) {
   Hclient := NewHttpClient()
   Convey("test sub unlock", t, func() {
       ctx, cancel := context.WithTimeout(context.Background(), waitForEventTimeout)
       defer cancel()

       eventCh, err := Hclient.Subscribe(ctx, "helper", "tm.event='Unlock'")
       if err != nil {
           t.Error("test sub unlock error", err)
       }
       go func() {
           stat, err := Hclient.Status()
           if err == nil {
               for {
                   select {
                   case evt := <-eventCh:
                       evData := evt.Data.(types.EventDataRoundState)
                       fmt.Printf("%d", evData.Height)
                       // So(evData.Height, ShouldBeGreaterThan, stat.SyncInfo.LatestBlockHeight)
                       return
                   }
               }
           } else {
               So(stat, ShouldBeError)
               return
           }
       }()
   })
   time.Sleep(10 *time.Second)
}

func TestSubEventValidatorSetUpdates(t *testing.T) {
   Hclient := NewHttpClient()
   Convey("test sub validator set updates", t, func() {
       ctx, cancel := context.WithTimeout(context.Background(), waitForEventTimeout)
       defer cancel()
       eventCh, err := Hclient.Subscribe(ctx, "helper", "tm.evnet='ValidatorSetUpdates'")
       if err != nil {
           t.Error("test sub validator updates error", err)
       }
       go func() {
           for {
               select {
               case evt := <-eventCh:
                   evData := evt.Data.(types.EventDataValidatorSetUpdates)
                   fmt.Printf("%d", len(evData.ValidatorUpdates))
                   //So(len(evData.ValidatorUpdates), ShouldNotBeEmpty)
                   return
               }
           }
       }()
   })
   time.Sleep(10 *time.Second)
}

func TestSubEventCompleteProposal(t *testing.T) {
   Hclient := NewHttpClient()
   Convey("test sub complete proposal", t, func() {
       ctx, cancel := context.WithTimeout(context.Background(), waitForEventTimeout)
       defer cancel()
       eventCh, err := Hclient.Subscribe(ctx, "helper", "tm.event='CompleteProposal'")
       if err != nil {
           t.Error("test  sub complete proposal error", err)
       }
       go func() {
           for {
               select {
               case evt := <-eventCh:
                   evData := evt.Data.(types.EventDataCompleteProposal)
                   fmt.Printf("%d", evData.Height)
                   //So(evData.Height, ShouldNotBeEmpty)
                   return
               }
           }
       }()
   })
   time.Sleep(10 *time.Second)
}

func TestSubEventEventNewRound(t *testing.T) {
   Hclient := NewHttpClient()
   Convey("event new round", t, func() {
       ctx, cancel := context.WithTimeout(context.Background(), waitForEventTimeout)
       defer cancel()
       eventCh, err := Hclient.Subscribe(ctx, "helper", "tm.event='NewRound'")
       if err != nil {
           t.Error("event new round", err)
       }
       go func() {
           for {
               select {
               case evt := <-eventCh:
                   evData := evt.Data.(types.EventDataNewRound)
                   fmt.Printf("%d", evData.Height)
                   //So(evData.Height, ShouldNotBeEmpty)
                   return
               }
           }
       }()
   })
   time.Sleep(10 *time.Second)
}

func TestSubEventEventNewRoundStep(t *testing.T) {
   Hclient := NewHttpClient()
   ctx, cancel := context.WithTimeout(context.Background(), waitForEventTimeout)
   defer cancel()
   eventCh, err := Hclient.Subscribe(ctx, "helper", "tm.event='NewRoundStep'")
   if err != nil {
       t.Error("event new round step", err)
   }
   go func() {
       for {
           select {
           case evt := <-eventCh:
               evData := evt.Data.(types.EventDataRoundState)
               fmt.Printf("%d", evData.Height)
               //So(evData.Height, ShouldNotBeEmpty)
               return
           }
       }
   }()
   time.Sleep(10 *time.Second)
}

func TestSubEventValidBlock(t *testing.T) {
   Hclient := NewHttpClient()
   ctx, cancel := context.WithTimeout(context.Background(), waitForEventTimeout)
   defer cancel()
   eventCh, err := Hclient.Subscribe(ctx, "helper", "tm.event='ValidBlock'")
   if err != nil {
       t.Error("event valid block", err)
   }
   go func() {
       for {
           select {
           case evt := <-eventCh:
               evData := evt.Data.(types.EventDataRoundState)
               fmt.Printf("%d", evData.Height)
               return
           }
       }
   }()
   time.Sleep(10 *time.Second)
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
   time.Sleep(10 *time.Second)
}

func TestSubEventPolka(t *testing.T) {
   Hclient := NewHttpClient()
   ctx, cancel := context.WithTimeout(context.Background(), waitForEventTimeout)
   defer cancel()
   eventCh, err := Hclient.Subscribe(ctx, "helper", "tm.event='Polka'")
   if err != nil {
       t.Error("event polka", err)
   }
   go func() {
       for {
           select {
           case evt := <-eventCh:
               evData := evt.Data.(types.EventDataRoundState)
               fmt.Printf("%d", evData.Height)
               return
           }
       }
   }()
   time.Sleep(10 *time.Second)
}

func TestSubEventTimeoutPropose(t *testing.T) {
   Hclient := NewHttpClient()
   ctx, cancel := context.WithTimeout(context.Background(), waitForEventTimeout)
   defer cancel()
   eventCh, err := Hclient.Subscribe(ctx, "helper", "tm.event='TimeoutPropose'")
   if err != nil {
       t.Error("event TimeoutPropose", err)
   }
   go func() {
       for {
           select {
           case evt := <-eventCh:
               evData := evt.Data.(types.EventDataRoundState)
               fmt.Printf("%d", evData.Height)
               err = Hclient.Unsubscribe(ctx, "helper", "tm.event='TimeoutPropose'")
               if err != nil {
                   t.Error("unsubscribe error")
               }
           }
       }
   }()
   time.Sleep(10 *time.Second)
}

func TestSubEventTimeoutWait(t *testing.T) {
   Hclient := NewHttpClient()
   ctx, cancel := context.WithTimeout(context.Background(), waitForEventTimeout)
   defer cancel()
   eventCh, err := Hclient.Subscribe(ctx, "helper", "tm.event='TimeoutWait'")
   if err != nil {
       t.Error("event TimeoutWait", err)
   }
   go func() {
       for {
           select {
           case evt := <-eventCh:
               evData := evt.Data.(types.EventDataRoundState)
               fmt.Printf("%d", evData.Height)
               return
           }
       }
   }()
   time.Sleep(10 * time.Second)
}