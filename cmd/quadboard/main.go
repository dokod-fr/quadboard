package main

import (
	"log/slog"
	"os"

	cli "github.com/dokod-fr/quadboard/internal/cli"
	"github.com/dokod-fr/quadboard/internal/config"
	"github.com/dokod-fr/quadboard/internal/logging"
)

func main() {
	os.Exit(run())
}

func run() int {
	logger := logging.New()
	slog.SetDefault(logger)

	_, err := config.Load()
	if err != nil {
		slog.Error("config load failed", "error", err)
		return 1
	}

	if err := cli.Execute(); err != nil {
		slog.Error("application terminated", "error", err)
		return 1
	}

	return 0
}
