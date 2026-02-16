package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/velocity-trinity/core/pkg/analyzer"
	"github.com/velocity-trinity/core/pkg/config"
	"github.com/velocity-trinity/core/pkg/logger"
)

var rootCmd = &cobra.Command{
	Use:   "dependency-ci",
	Short: "Smart Dependency Analyzer for CI Pipelines",
	Long:  `Analyzes your codebase to determine which tests need to run based on changed files.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var analyzeCmd = &cobra.Command{
	Use:   "analyze [file]",
	Short: "Analyze a single file for dependencies",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		logger.Log.Info("Analyzing file: " + filePath)

		deps, err := analyzer.AnalyzeFile(filePath)
		if err != nil {
			logger.Log.Error("Analysis failed: " + err.Error())
			os.Exit(1)
		}

		if len(deps) == 0 {
			fmt.Println("No relative dependencies found.")
		} else {
			fmt.Printf("Dependencies for %s:\n", filePath)
			for _, dep := range deps {
				fmt.Println(" - " + dep)
			}
		}
	},
}

func main() {
	// Initialize Config & Logger
	cfg, err := config.Load("dependency-ci")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Config error: %v\n", err)
		// Don't exit, use defaults
	}

	env := "development"
	if cfg != nil {
		env = cfg.Env
	}
	logger.Init(env)
	defer logger.Sync()

	rootCmd.AddCommand(analyzeCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
