package transport

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/rpc"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

	// Post-Sync Command Execution
	if req.PostSyncCommand != "" {
		logger.Log.Info("Executing post-sync command: " + req.PostSyncCommand)
		
		// Split command string carefully (this is simple, ignores quotes)
		parts := strings.Fields(req.PostSyncCommand)
		head := parts[0]
		args := parts[1:]

		cmd := exec.Command(head, args...)
		cmd.Dir = s.BasePath // Execute in project root
		
		output, err := cmd.CombinedOutput()
		if err != nil {
			logger.Log.Error("Post-sync command failed: " + err.Error())
			resp.Success = false
			resp.Message = fmt.Sprintf("File synced, but command failed: %s\nOutput: %s", err.Error(), string(output))
			return nil // Return nil error so RPC succeeds, but resp indicates failure
		}
		
		logger.Log.Info("Command executed successfully")
		resp.Message = fmt.Sprintf("File synced and command executed: %s", string(output))
	} else {
		resp.Message = "File synced successfully"
	}

	resp.Success = true
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
