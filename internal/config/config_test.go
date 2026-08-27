package config

import (
	"log/slog"
	"testing"
	"time"
)

func validRawConfig() rawConfig {
	var raw rawConfig
	raw.Network.Address = "127.0.0.1:3223"
	raw.Network.MaxConnections = 100
	raw.Network.MaxMessageSize = "4KB"
	raw.Network.IdleTimeout = "30s"
	raw.Logging.Level = "info"
	return raw
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(raw *rawConfig)
		want    Config
		wantErr bool
	}{
		{
			name:   "valid config",
			mutate: func(raw *rawConfig) {},
			want: Config{
				Network: NetworkConfig{
					Address:        "127.0.0.1:3223",
					MaxConnections: 100,
					MaxMessageSize: 4 * 1024,
					IdleTimeout:    30 * time.Second,
				},
				Logging: LoggingConfig{
					Level: slog.LevelInfo,
				},
			},
		},
		{
			name: "invalid max_message_size",
			mutate: func(raw *rawConfig) {
				raw.Network.MaxMessageSize = "not-a-size"
			},
			wantErr: true,
		},
		{
			name: "invalid idle_timeout",
			mutate: func(raw *rawConfig) {
				raw.Network.IdleTimeout = "not-a-duration"
			},
			wantErr: true,
		},
		{
			name: "invalid logging level",
			mutate: func(raw *rawConfig) {
				raw.Logging.Level = "not-a-level"
			},
			wantErr: true,
		},
		{
			name: "empty raw config",
			mutate: func(raw *rawConfig) {
				*raw = rawConfig{}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := validRawConfig()
			tt.mutate(&raw)

			got, err := normalize(raw)

			if tt.wantErr {
				if err == nil {
					t.Fatal("normalize() error = nil, want error")
				}
				if got != nil {
					t.Error("normalize() config != nil on error")
				}
				return
			}

			if err != nil {
				t.Fatalf("normalize() unexpected error = %v", err)
			}
			if got == nil {
				t.Fatal("normalize() config = nil, want non-nil")
			}
			if *got != tt.want {
				t.Errorf("normalize() = %+v, want %+v", *got, tt.want)
			}
		})
	}
}
