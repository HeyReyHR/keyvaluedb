package database

import (
	"errors"
	"fmt"
	"log/slog"

	parser2 "github.com/HeyReyHR/keyvaluedb/internal/database/compute/parser"
	"github.com/HeyReyHR/keyvaluedb/internal/database/storage"
)

type computeLayer interface {
	Parse(string) (parser2.Query, error)
}

type storageLayer interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Del(key string) error
}

type Database struct {
	computeLayer computeLayer
	storageLayer storageLayer
	logger       *slog.Logger
}

func NewDatabase(computeLayer computeLayer, storageLayer storageLayer, logger *slog.Logger) (*Database, error) {
	if computeLayer == nil {
		return nil, errors.New("invalid compute layer")
	}

	if storageLayer == nil {
		return nil, errors.New("invalid storage layer")
	}

	if logger == nil {
		logger = slog.Default()
	}
	return &Database{
		computeLayer: computeLayer,
		storageLayer: storageLayer,
		logger:       logger,
	}, nil
}

func (d *Database) HandleQuery(input string) string {
	query, err := d.computeLayer.Parse(input)
	if err != nil {
		return fmt.Sprintf("[error] %s", err) // Почему err.Error() а не просто err?
	}

	switch query.CommandId() {
	case parser2.GetCommandId:
		return d.handleGetQuery(query)
	case parser2.DelCommandId:
		return d.handleDelQuery(query)
	case parser2.SetCommandId:
		return d.handleSetQuery(query)
	default:
		d.logger.Error("compute layer is incorrect", slog.Int("command_id", query.CommandId()))
		return "ERROR internal error"
	}
}

func (d *Database) handleGetQuery(query parser2.Query) string {
	args := query.Arguments()
	value, err := d.storageLayer.Get(args[0])
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return fmt.Sprintf("NOT FOUND")
		}
		return fmt.Sprintf("ERROR %s", err)
	}

	return fmt.Sprintf("OK %s", value)
}

func (d *Database) handleSetQuery(query parser2.Query) string {
	args := query.Arguments()

	if err := d.storageLayer.Set(args[0], args[1]); err != nil {
		return fmt.Sprintf("ERROR %s", err)
	}

	return fmt.Sprintf("OK")
}

func (d *Database) handleDelQuery(query parser2.Query) string {
	args := query.Arguments()

	if err := d.storageLayer.Del(args[0]); err != nil {
		return fmt.Sprintf("ERROR %s", err)
	}

	return fmt.Sprintf("OK")
}
