package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

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

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run tests only for changed files",
	Long: `Example: dependency-ci run --cmd="npm test" --files="src/foo.ts src/bar.ts"`,
	Run: func(cmd *cobra.Command, args []string) {
		filesStr, _ := cmd.Flags().GetString("files")
		testCmd, _ := cmd.Flags().GetString("cmd")

		if filesStr == "" {
			logger.Log.Info("No files changed. Skipping tests.")
			return
		}

		files := strings.Split(filesStr, " ")
		testFiles, err := analyzer.FindTestFiles(files)
		if err != nil {
			logger.Log.Fatal("Failed to find tests: " + err.Error())
		}

		if len(testFiles) == 0 {
			logger.Log.Info("No tests affected by these changes.")
			return
		}

		logger.Log.Info(fmt.Sprintf("Running %d test files...", len(testFiles)))
		
		// Construct the command: npm test file1 file2 ...
		fullCmd := fmt.Sprintf("%s %s", testCmd, strings.Join(testFiles, " "))
		logger.Log.Info("Executing: " + fullCmd)

		// execute command
		parts := strings.Fields(fullCmd)
		head := parts[0]
		parts = parts[1:]

		c := exec.Command(head, parts...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		
		if err := c.Run(); err != nil {
			logger.Log.Error("Tests failed!")
			os.Exit(1)
		}
		
		logger.Log.Info("All relevant tests passed!")
	},
}

func init() {
	runCmd.Flags().String("files", "", "Space-separated list of changed files")
	runCmd.Flags().String("cmd", "npm test", "Base test command (e.g., 'npm test', 'pytest')")
}

func main() {
	// Initialize Config & Logger
	cfg, err := config.Load("dependency-ci")
	if err != nil {
		// Non-fatal config error
	}

	env := "development"
	if cfg != nil {
		env = cfg.Env
	}
	logger.Init(env)
	defer logger.Sync()

	rootCmd.AddCommand(analyzeCmd)
	rootCmd.AddCommand(runCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
