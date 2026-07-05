package main

import (
	"log/slog"
	"os"

	"github.com/dokod-fr/quadboard/internal/cli"
)

func main() {
	os.Exit(run())
}

func run() int {

	if err := cli.Execute(); err != nil {
		slog.Error("command failed", "error", err)
		return 1
	}
	return 0
}
