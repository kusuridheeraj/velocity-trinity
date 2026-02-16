package transport

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/rpc"
	"os"
	"path/filepath"

	"github.com/velocity-trinity/core/pkg/logger"
)

// LivePatchServer handles RPC calls
type LivePatchServer struct {
	BasePath string
}

// SyncFile is the RPC method called by the client
func (s *LivePatchServer) SyncFile(req *FileSyncRequest, resp *FileSyncResponse) error {
	fullPath := filepath.Join(s.BasePath, req.RelativePath)
	
	// Security check: Prevent directory traversal outside base path
	if !filepath.IsAbs(fullPath) {
		absBase, _ := filepath.Abs(s.BasePath)
		absFull, _ := filepath.Abs(fullPath)
		if len(absFull) < len(absBase) || absFull[:len(absBase)] != absBase {
			return fmt.Errorf("security violation: path traversal detected")
		}
	}

	logger.Log.Info("Syncing file: " + fullPath)
	
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		resp.Success = false
		resp.Message = "Failed to create directory: " + err.Error()
		return err
	}

	if err := ioutil.WriteFile(fullPath, req.Content, os.FileMode(req.Mode)); err != nil {
		resp.Success = false
		resp.Message = "Failed to write file: " + err.Error()
		return err
	}

	resp.Success = true
	resp.Message = "File synced successfully"
	return nil
}

// StartServer starts the RPC server
func StartServer(port string, basePath string, tlsConfig *tls.Config) error {
	server := &LivePatchServer{BasePath: basePath}
	rpc.Register(server)

	listener, err := tls.Listen("tcp", ":"+port, tlsConfig)
	if err != nil {
		return err
	}
	
	logger.Log.Info("LivePatch Agent listening on port " + port)
	
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Log.Error("Accept error: " + err.Error())
			continue
		}
		go rpc.ServeConn(conn)
	}
}
