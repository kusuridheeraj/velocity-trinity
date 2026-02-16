package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log is the global logger instance
var Log *zap.Logger

// Init initializes the logger based on the environment
func Init(env string) {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Scrub sensitive data (rudimentary implementation for now)
	// In a real implementation, this would involve a custom zapcore.Core
	// wrapping the standard one.
	
	var err error
	Log, err = config.Build()
	if err != nil {
		os.Stderr.WriteString("Failed to initialize logger: " + err.Error())
		os.Exit(1)
	}
	
	zap.ReplaceGlobals(Log)
}

// Sync flushes the logger
func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
