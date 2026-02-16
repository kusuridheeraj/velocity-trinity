package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/velocity-trinity/core/pkg/config"
	"github.com/velocity-trinity/core/pkg/logger"
	"github.com/velocity-trinity/core/pkg/scheduler"
	"github.com/velocity-trinity/core/pkg/webhook"
)

var rootCmd = &cobra.Command{
	Use:   "quantum-merge",
	Short: "Speculative Merge Queue Manager",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Quantum Merge Scheduler",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Log.Info("Starting Quantum Merge Scheduler...")
		
		queue := scheduler.NewMemoryQueue(100)
		
		// Start Workers
		for i := 0; i < 4; i++ {
			go scheduler.RunWorker(i, queue)
		}
		
		// Start Webhook Server
		server := webhook.NewServer(queue, "8090")
		if err := server.ListenAndServe(); err != nil {
			logger.Log.Fatal("Webhook server crashed: " + err.Error())
		}
	},
}

func main() {
	// Initialize Config & Logger
	cfg, _ := config.Load("quantum-merge")
	env := "development"
	if cfg != nil {
		env = cfg.Env
	}
	logger.Init(env)
	defer logger.Sync()

	rootCmd.AddCommand(serveCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
