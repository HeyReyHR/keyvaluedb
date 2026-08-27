package initialization

import (
	"log/slog"
	"os"

	"github.com/HeyReyHR/keyvaluedb/internal/config"
)

func CreateLogger(cfg config.LoggingConfig) *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: cfg.Level,
	}))
}
