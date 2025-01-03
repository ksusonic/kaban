package logger

import (
	"io"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func New(isDebug bool) *slog.Logger {
	level := slog.LevelInfo
	if isDebug {
		level = slog.LevelDebug
	}

	return slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: level}))
}

func NewDisabled() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
