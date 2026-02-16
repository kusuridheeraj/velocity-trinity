package transport

import "time"

// FileSyncRequest represents a request to sync a file
type FileSyncRequest struct {
	RelativePath string
	Content      []byte
	Mode         uint32
	Timestamp    time.Time
	
	// PostSyncCommand: Command to execute after file sync
	PostSyncCommand string
}

// FileSyncResponse represents the result of a sync operation
type FileSyncResponse struct {
	Success bool
	Message string
}

// CommandRequest represents a request to run a command (e.g. restart server)
type CommandRequest struct {
	Command []string
	Timeout int // seconds
}

// CommandResponse represents the result of a command execution
type CommandResponse struct {
	Success  bool
	ExitCode int
	Output   string
	Error    string
}
