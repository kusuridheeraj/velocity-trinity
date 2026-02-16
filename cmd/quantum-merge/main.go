package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/velocity-trinity/core/pkg/config"
	"github.com/velocity-trinity/core/pkg/logger"
	"github.com/velocity-trinity/core/pkg/scheduler"
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
		scheduler.Run(4) // Start 4 workers
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
