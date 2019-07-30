package tester

import (
    "fmt"
    . "github.com/smartystreets/goconvey/convey"
    "github.com/tendermint/tendermint/rpc/client"
    "testing"
)

func NewHttpClient() *client.HTTP {
    httpClient := client.NewHTTP(node1+":26657", "/websocket")
    err := httpClient.Start()
    if err != nil {
        fmt.Printf("new http error %v", err)
    }
    return httpClient
}

func TestStatus(t *testing.T) {
    Hclient := NewHttpClient()
    Convey("test status", t, func() {
        stat, err := Hclient.Status()
        if err == nil {
            fmt.Printf("stat ok! node Height=%d", stat.SyncInfo.LatestBlockHash)
        } else {
            t.Error("status error", err)
        }
    })
}

func TestBlockResults(t *testing.T) {
    Hclient := NewHttpClient()
    Convey("test transaction tx", t, func() {
        stat, err := Hclient.Status()
        if err == nil {
            height := stat.SyncInfo.LatestBlockHeight
            if height > 1 {
                height = height - 1
            }
            result, err := Hclient.BlockResults(&height)
            if err == nil {
                if result.Height == height {
                    // t.Log("BlockResults Ok")
                }
            } else {
                t.Error("BlockResults error", err)
            }
        }
    })
}

func TestBlockAndTx(t *testing.T) {
    Hclient := NewHttpClient()
    Convey("test transaction tx and block", t, func() {
        stat, err := Hclient.Status()
        if err == nil {
            height := stat.SyncInfo.LatestBlockHeight
            if height > 1 {
                height = height - 1
            }
            result, err := Hclient.Block(&height)
            if err != nil {
                t.Error("block error ", err)
            } else {
                if result.Block.Height == height {
                    t.Log("block OK")
                }
                for _, tx := range result.Block.Txs {
                    res, err := Hclient.Tx(tx.Hash(), false)
                    if err == nil {
                        So(res.Hash.String(), ShouldEqual, fmt.Sprintf("%X", tx.Hash()))
                    } else {
                        t.Error("Tx error", err)
                    }
                }
            }

        }
    })
}

func TestSearch(t *testing.T) {
    Hclient := NewHttpClient()
    Convey("test transaction tx and block", t, func() {
        stat, err := Hclient.Status()
        if err != nil {
            t.Error("client status ", err)
        } else {
            queryStr := fmt.Sprintf("tx.height>=%d AND tx.height <=%d", stat.SyncInfo.LatestBlockHeight, stat.SyncInfo.LatestBlockHeight)
            res, err := Hclient.TxSearch(queryStr, false, 1, 30)
            if err != nil {
                t.Error("res txsearch ", err)
            } else {
                fmt.Sprintf("%d", res.TotalCount)
            }
        }
    })
}

func TestCommit(t *testing.T) {
    Hclient := NewHttpClient()
    Convey("test commit ", t, func() {
        stat, err := Hclient.Status()
        if err != nil {
            t.Error("client status ", err)
        } else {
            res, err := Hclient.Commit(&stat.SyncInfo.LatestBlockHeight)
            if err != nil {
                t.Error("test commit error", err)
            } else {
                t.Log(res.Height)
            }
        }
    })
}

func TestBlockchainInfo(t *testing.T) {
    Hclient := NewHttpClient()
    Convey("test sub Height", t, func() {
        stat, err := Hclient.Status()
        if err != nil {
            t.Error("status", err)
        } else {
            height := stat.SyncInfo.LatestBlockHeight
            var min int64
            if height > 10 {
                min = stat.SyncInfo.LatestBlockHeight - 10
            }
            res, err := Hclient.BlockchainInfo(min, stat.SyncInfo.LatestBlockHeight)
            if err == nil {
                So(len(res.BlockMetas), ShouldEqual, 11)
            } else {
                t.Error("BlockchainInfo err", err)
            }
        }
    })
}

func TestGenesis(t *testing.T) {
    Hclient := NewHttpClient()
    Convey("test transaction tx and block", t, func() {
        res, err := Hclient.Genesis()
        if err != nil {
            t.Error("client Genesis", err)
        } else {
            t.Log(res.Genesis.ChainID)
        }
    })
}

func TestHealth(t *testing.T) {
    Hclient := NewHttpClient()
    Convey("test transaction tx and block", t, func() {
        _, err := Hclient.Health()
        if err != nil {
            t.Error("client Genesis", err)
        } else {
            t.Log("health ok")
        }
    })
}

func TestBroadcastTxAsync(t *testing.T) {
    Hclient := NewHttpClient()
    Convey("test transaction tx and block", t, func() {
        stat, err := Hclient.Status()
        if err == nil {
            height := stat.SyncInfo.LatestBlockHeight
            if height > 1 {
                height = height - 1
            }
            result, err := Hclient.Block(&height)
            if err != nil {
                t.Error("block error ", err)
            } else {
                if result.Block.Height == height {
                    t.Log("block OK")
                }
                for _, tx := range result.Block.Txs {
                    _, err := Hclient.BroadcastTxAsync(tx)
                    So(err, ShouldBeError)
                }
            }

        }
    })
}

func TestBroadcastTxCommit(t *testing.T) {
    Hclient := NewHttpClient()
    Convey("test transaction tx and block", t, func() {
        stat, err := Hclient.Status()
        if err == nil {
            height := stat.SyncInfo.LatestBlockHeight
            if height > 1 {
                height = height - 1
            }
            result, err := Hclient.Block(&height)
            if err != nil {
                t.Error("block error ", err)
            } else {
                So(result.Block.Height, ShouldEqual, height)
                for _, tx := range result.Block.Txs {
                    _, err := Hclient.BroadcastTxCommit(tx)
                    So(err, ShouldBeError)
                }
            }
        }
    })
}

func TestBroadcastTxSync(t *testing.T) {
    Hclient := NewHttpClient()
    Convey("test transaction tx and block", t, func() {
        stat, err := Hclient.Status()
        if err == nil {
            height := stat.SyncInfo.LatestBlockHeight
            if height > 1 {
                height = height - 1
            }
            result, err := Hclient.Block(&height)
            if err != nil {
                t.Error("block error ", err)
            } else {
                So(result.Block.Height, ShouldEqual, height)
                for _, tx := range result.Block.Txs {
                    _, err := Hclient.BroadcastTxSync(tx)
                    So(err, ShouldBeError)
                }
            }
        }
    })
}
