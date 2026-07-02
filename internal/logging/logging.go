package logging

import (
	"os"

	"log/slog"
)

func New() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}
