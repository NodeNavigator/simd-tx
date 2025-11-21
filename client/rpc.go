package client

import (
	"context"
	"fmt"
	"time"

	rpc "github.com/cometbft/cometbft/rpc/client/http"
	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
)

type RPCClient struct {
	client *rpc.HTTP
}

// ConnectRPCWithTimeout connects RPC client connection with timeout.
func ConnectRPCWithTimeout(addr string, timeout time.Duration) (*RPCClient, error) {
	rpcClient, err := rpc.NewWithTimeout(addr, "/websocket", uint(timeout.Seconds()))
	if err != nil {
		return nil, err
	}
	return &RPCClient{client: rpcClient}, nil
}

// Block returns block information for the height.
func (c *RPCClient) Block(ctx context.Context, height int64) (*tmctypes.ResultBlock, error) {
	return c.client.Block(ctx, &height)
}

// LatestBlockHeight returns the latest block height on the network.
func (c *RPCClient) LatestBlockHeight(ctx context.Context) (int64, error) {
	resp, err := c.client.Status(ctx)
	if err != nil {
		return 0, err
	}

	return resp.SyncInfo.LatestBlockHeight, nil
}

// NetworkChainID returns network chain id.
func (c *RPCClient) NetworkChainID(ctx context.Context) (string, error) {
	status, err := c.client.Status(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get status: %w", err)
	}

	return status.NodeInfo.Network, nil
}

// Status returns the status of the blockchain network.
func (c *RPCClient) Status(ctx context.Context) (*tmctypes.ResultStatus, error) {
	return c.client.Status(ctx)
}
