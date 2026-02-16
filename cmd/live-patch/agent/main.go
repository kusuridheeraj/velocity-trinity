package main

import (
	"crypto/tls"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/velocity-trinity/core/pkg/config"
	"github.com/velocity-trinity/core/pkg/logger"
	"github.com/velocity-trinity/core/pkg/transport"
	"github.com/velocity-trinity/core/pkg/utils"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "live-patch-agent",
		Short: "Agent for receiving hot patches inside containers",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Log.Info("Starting LivePatch Agent...")

			// Development Mode: Generate certificates if they don't exist
			if _, err := os.Stat("server.crt"); os.IsNotExist(err) {
				logger.Log.Info("Generating self-signed certificates for development...")
				if err := utils.GenerateSelfSignedCert("server.crt", "server.key"); err != nil {
					// We can't log fatal here because logger is initialized later? No, it's global.
					// But for simplicity let's use fmt if logger fails
					fmt.Println("Failed to generate certs: " + err.Error())
					os.Exit(1)
				}
			}

			// Load TLS Config
			cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
			if err != nil {
				logger.Log.Fatal("Failed to load key pair: " + err.Error())
			}
			tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}

			// Start Server
			port := os.Getenv("PORT")
			if port == "" {
				port = "8080"
			}
			
			basePath := os.Getenv("BASE_PATH")
			if basePath == "" {
				basePath = "."
			}

			if err := transport.StartServer(port, basePath, tlsConfig); err != nil {
				logger.Log.Fatal("Server crashed: " + err.Error())
			}
		},
	}

	// Initialize Config & Logger
	cfg, _ := config.Load("live-patch-agent")
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
