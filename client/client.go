package client

import (
	"context"
	"fmt"
	"time"

	"github.com/Ankr-network/ankr-chain/common/code"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/pubsub/query"
	"github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
)

type Client struct {
	cHttp *client.HTTP
	cdc   *amino.Codec
}

func NewClient(nodeUrl string) *Client {
	cHttp := client.NewHTTP(nodeUrl, "/websocket")

	return &Client{cHttp, amino.NewCodec()}
}

func (c *Client) Query(path string, req interface{}, resp interface{}) (err error) {
	reqDataBytes, err := c.cdc.MarshalJSON(req)
	if err != nil {
		return err
	}

	resultQ, err := c.cHttp.ABCIQuery(path, reqDataBytes)
	if err != nil {
		return err
	}

	if resultQ.Response.Code != code.CodeTypeOK {
		return fmt.Errorf("Client query response code not ok, code=%d, log=%s", resultQ.Response.Code, resultQ.Response.Log)
	}

	c.cdc.UnmarshalJSON(resultQ.Response.Value, resp)

	return nil
}

func (c *Client) BroadcastTxCommit(txBytes []byte) (txHash string, commitHeight int64, log string, err error) {
	result, err := c.cHttp.BroadcastTxCommit(txBytes)
	if err != nil {
		return "", -1, "", err
	}

	if result.CheckTx.Code != code.CodeTypeOK {
		return "", -1, "", fmt.Errorf("Client BroadcastTxCommit CheckTx response code not ok, code=%d, log=%s", result.CheckTx.Code, result.CheckTx.Log)
	}

	if result.DeliverTx.Code != code.CodeTypeOK {
		return "", -1, "", fmt.Errorf("Client BroadcastTxCommit DeliverTx response code not ok, code=%d, log=%s", result.DeliverTx.Code, result.DeliverTx.Log)
	}

	err = client.WaitForHeight(c.cHttp, result.Height+1, nil)
	if err != nil {
		return result.Hash.String(), result.Height, "", err
	}

	return result.Hash.String(), result.Height+1, result.DeliverTx.Log, nil
}

func (c *Client) Status() (*ctypes.ResultStatus, error) {
	return c.cHttp.Status()
}

func (c *Client) ABCIInfo() (*ctypes.ResultABCIInfo, error) {
	return c.cHttp.ABCIInfo()
}

func (c *Client) BroadcastTxAsync(tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	return c.cHttp.BroadcastTxAsync(tx)
}

func (c *Client) BroadcastTxSync(tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	return c.cHttp.BroadcastTxSync(tx)
}

func (c *Client) UnconfirmedTxs(limit int) (*ctypes.ResultUnconfirmedTxs, error) {
	return c.cHttp.UnconfirmedTxs(limit)
}

func (c *Client) NumUnconfirmedTxs() (*ctypes.ResultUnconfirmedTxs, error) {
	return c.cHttp.NumUnconfirmedTxs()
}

func (c *Client) NetInfo() (*ctypes.ResultNetInfo, error) {
	return c.cHttp.NetInfo()
}

func (c *Client) DumpConsensusState() (*ctypes.ResultDumpConsensusState, error) {
	return c.cHttp.DumpConsensusState()
}

func (c *Client) ConsensusState() (*ctypes.ResultConsensusState, error) {
	return c.cHttp.ConsensusState()
}

func (c *Client) Health() (*ctypes.ResultHealth, error) {
	return c.cHttp.Health()
}

func (c *Client) BlockchainInfo(minHeight, maxHeight int64) (*ctypes.ResultBlockchainInfo, error) {
	return c.cHttp.BlockchainInfo(minHeight, maxHeight)
}

func (c *Client) Genesis() (*ctypes.ResultGenesis, error) {
	return c.cHttp.Genesis()
}

func (c *Client) Block(height *int64) (*ctypes.ResultBlock, error) {
	return c.cHttp.Block(height)
}

func (c *Client) BlockResults(height *int64) (*ctypes.ResultBlockResults, error) {
	return c.cHttp.BlockResults(height)
}

func (c *Client) Commit(height *int64) (*ctypes.ResultCommit, error) {
	return c.cHttp.Commit(height)
}

func (c *Client) Tx(hash []byte, prove bool) (*ctypes.ResultTx, error) {
	return c.cHttp.Tx(hash, prove)
}

func (c *Client) TxSearch(query string, prove bool, page, perPage int) (*ctypes.ResultTxSearch, error) {
	return c.cHttp.TxSearch(query, prove, page, perPage)
}

func (c *Client) Validators(height *int64) (*ctypes.ResultValidators, error) {
	return c.cHttp.Validators(height)
}

func (c *Client) SubscribeAndWait(subscriber string, queryStr string, waitTimeOut time.Duration, maxSubCap int, out chan ctypes.ResultEvent) error {
	c.cHttp.Start()

	q := query.MustParse(queryStr)

	ctx, cancel := context.WithTimeout(context.Background(), waitTimeOut)
	defer cancel()

	eventCh, err := c.cHttp.Subscribe(ctx, subscriber, q.String(), maxSubCap)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case evt := <-eventCh:
				out <- evt
			}
		}
	}()

	return nil
}