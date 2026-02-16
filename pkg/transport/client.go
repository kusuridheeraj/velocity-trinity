package transport

import (
	"crypto/tls"
	"fmt"
	"net/rpc"
	"time"

	"github.com/velocity-trinity/core/pkg/logger"
)

// LivePatchClient handles RPC connections to the agent
type LivePatchClient struct {
	client *rpc.Client
}

// NewClient creates a new LivePatchClient
func NewClient(addr string, tlsConfig *tls.Config) (*LivePatchClient, error) {
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("connection error: %v", err)
	}

	client := rpc.NewClient(conn)
	return &LivePatchClient{client: client}, nil
}

// SyncFile syncs a local file to the remote agent
func (c *LivePatchClient) SyncFile(req *FileSyncRequest) (*FileSyncResponse, error) {
	var resp FileSyncResponse
	
	start := time.Now()
	err := c.client.Call("LivePatchServer.SyncFile", req, &resp)
	duration := time.Since(start)

	if err != nil {
		logger.Log.Error("RPC error: " + err.Error())
		return nil, err
	}

	logger.Log.Info(fmt.Sprintf("Synced %s in %v", req.RelativePath, duration))
	return &resp, nil
}

// Close closes the client connection
func (c *LivePatchClient) Close() error {
	return c.client.Close()
}
