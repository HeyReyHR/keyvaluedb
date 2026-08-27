package initialization

import (
	"errors"
	"log/slog"

	"github.com/HeyReyHR/keyvaluedb/internal/config"
	"github.com/HeyReyHR/keyvaluedb/internal/network"
)

const defaultServerAddress = ":3223"

func CreateNetwork(cfg *config.NetworkConfig, logger *slog.Logger) (*network.TCPServer, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	address := defaultServerAddress
	var options []network.TCPServerOption

	if cfg != nil {
		if cfg.Address != "" {
			address = cfg.Address
		}
		if cfg.MaxConnections != 0 {
			options = append(options, network.WithServerMaxConnections(uint(cfg.MaxConnections)))
		}
		if cfg.MaxMessageSize != 0 {
			options = append(options, network.WithServerBufferSize(uint(cfg.MaxMessageSize)))
		}
		if cfg.IdleTimeout != 0 {
			options = append(options, network.WithServerIdleTimeout(cfg.IdleTimeout))
		}
	}

	return network.NewTCPServer(address, logger, options...)
}
