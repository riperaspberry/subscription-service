package logger

import (
	"log/slog"
	"os"
)

func Init() *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	l := slog.New(handler)
	slog.SetDefault(l)

	return l
}
