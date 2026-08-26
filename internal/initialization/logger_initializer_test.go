package initialization

import (
	"context"
	"log/slog"
	"testing"

	"github.com/HeyReyHR/keyvaluedb/internal/config"
)

func TestCreateLogger(t *testing.T) {
	tests := []struct {
		name        string
		level       slog.Level
		checkLevel  slog.Level
		wantEnabled bool
	}{
		{
			name:        "info level enabled for info",
			level:       slog.LevelInfo,
			checkLevel:  slog.LevelInfo,
			wantEnabled: true,
		},
		{
			name:        "warn level disables debug",
			level:       slog.LevelWarn,
			checkLevel:  slog.LevelDebug,
			wantEnabled: false,
		},
		{
			name:        "warn level enables error",
			level:       slog.LevelWarn,
			checkLevel:  slog.LevelError,
			wantEnabled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := CreateLogger(config.LoggingConfig{Level: tt.level})

			if logger == nil {
				t.Fatal("CreateLogger() = nil, want non-nil logger")
			}

			got := logger.Enabled(context.Background(), tt.checkLevel)
			if got != tt.wantEnabled {
				t.Errorf("Enabled(%v) = %v, want %v", tt.checkLevel, got, tt.wantEnabled)
			}
		})
	}
}
