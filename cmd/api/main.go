package main

import (
	"fmt"
	"log/slog"
	"logger-test/internal/server"
)

func main() {

	loggerConfig := server.LoggerConfig{
		LoggerType: "json",
		Level:      slog.LevelDebug,
		AddSource:  true,
	}

	logger, err := server.NewLogger(loggerConfig)
	if err != nil {
		panic(fmt.Sprintf("cannot create logger: %s", err))
	}

	logger.Info("starting server")
	logger.Debug("debugging server")
	logger.Warn("warning server")
	logger.Error("error server")
}
