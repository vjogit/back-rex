package utils

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log est l'instance globale de notre logger zap.
var Log *zap.Logger

func init() {
	logLevel := zapcore.InfoLevel
	if os.Getenv("DEBUG") != "" {
		logLevel = zapcore.DebugLevel
	}

	config := zap.NewProductionConfig()
	config.Level.SetLevel(logLevel)

	var err error
	Log, err = config.Build()
	if err != nil {
		log.Fatalf("failed to initialize zap logger: %v", err)
	}
}

// Sync force l'écriture de tous les logs en buffer.
func Sync() error {
	if Log != nil {
		return Log.Sync()
	}
	return nil
}
