package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

func setupLogger(verbose bool) *slog.Logger {
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}

	opts := &tint.Options{
		Level:      level,
		TimeFormat: time.Kitchen,
	}

	return slog.New(tint.NewHandler(os.Stderr, opts))
}
