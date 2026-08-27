package initialization

import (
	"context"
	"errors"
	"log/slog"

	"github.com/HeyReyHR/keyvaluedb/internal/config"
	"github.com/HeyReyHR/keyvaluedb/internal/database"
	"github.com/HeyReyHR/keyvaluedb/internal/database/compute/parser"
	"github.com/HeyReyHR/keyvaluedb/internal/database/storage"
	"github.com/HeyReyHR/keyvaluedb/internal/database/storage/engine"
	"github.com/HeyReyHR/keyvaluedb/internal/network"
)

type Initializer struct {
	server *network.TCPServer
	logger *slog.Logger
}

func NewInitializer(cfg *config.Config) (*Initializer, error) {
	if cfg == nil {
		return nil, errors.New("invalid config")
	}
	logger := CreateLogger(cfg.Logging)

	server, err := CreateNetwork(&cfg.Network, logger)
	if err != nil {
		return nil, err
	}

	initializer := &Initializer{
		logger: logger,
		server: server,
	}

	return initializer, nil
}

func (i *Initializer) StartDatabase(ctx context.Context) error {
	p := parser.NewParser(i.logger)
	e := engine.NewEngine()
	s := storage.NewStorage(e)

	db, err := database.NewDatabase(p, s, i.logger)
	if err != nil {
		return err
	}

	i.server.HandleQueries(ctx, func(ctx context.Context, query []byte) []byte {
		response := db.HandleQuery(string(query))
		return []byte(response)
	})

	return nil
}
