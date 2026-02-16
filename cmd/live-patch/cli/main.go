package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/velocity-trinity/core/pkg/config"
	"github.com/velocity-trinity/core/pkg/logger"
	"github.com/velocity-trinity/core/pkg/transport"
)

var targetAddr string

var rootCmd = &cobra.Command{
	Use:   "live-patch",
	Short: "LivePatch CLI Tool",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var syncCmd = &cobra.Command{
	Use:   "sync [file]",
	Short: "Sync a local file to the remote container",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		
		// Read file content
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			logger.Log.Fatal("Failed to read file: " + err.Error())
		}

		info, err := os.Stat(filePath)
		if err != nil {
			logger.Log.Fatal("Failed to stat file: " + err.Error())
		}

		// Connect to Agent
		// In production, use proper CA verification. For dev, skip verify.
		tlsConfig := &tls.Config{InsecureSkipVerify: true}
		
		client, err := transport.NewClient(targetAddr, tlsConfig)
		if err != nil {
			logger.Log.Fatal("Failed to connect to agent: " + err.Error())
		}
		defer client.Close()

		// Send Request
		req := &transport.FileSyncRequest{
			RelativePath: filePath, // Should be relative to project root
			Content:      content,
			Mode:         uint32(info.Mode()),
			Timestamp:    info.ModTime(),
		}

		resp, err := client.SyncFile(req)
		if err != nil {
			logger.Log.Fatal("Sync failed: " + err.Error())
		}

		if resp.Success {
			fmt.Printf("Successfully synced %s to %s\n", filePath, targetAddr)
		} else {
			fmt.Printf("Sync failed: %s\n", resp.Message)
		}
	},
}

func main() {
	syncCmd.Flags().StringVarP(&targetAddr, "target", "t", "localhost:8080", "Address of the LivePatch Agent")

	rootCmd.AddCommand(syncCmd)

	// Initialize Config & Logger
	cfg, _ := config.Load("live-patch")
	env := "development"
	if cfg != nil {
		env = cfg.Env
	}
	logger.Init(env)
	defer logger.Sync()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
