package main

import (
	"log/slog"
	"os"

	cli "github.com/dokod-fr/quadboard/internal/cli"
	"github.com/dokod-fr/quadboard/internal/logging"
)

func main() {
	os.Exit(run())
}

func run() int {
	logger := logging.New()
	slog.SetDefault(logger)

	if err := cli.Execute(); err != nil {
		slog.Error("application terminated", "error", err)
		return 1
	}

	return 0
}
