package parser

import (
	"errors"
	"log/slog"
	"strings"
)

var (
	ErrInvalidQuery      = errors.New("invalid query")
	ErrInvalidCommand    = errors.New("invalid command")
	ErrInvalidArgsNumber = errors.New("invalid arguments number")
)

type Parser struct {
	logger *slog.Logger
}

func NewParser(logger *slog.Logger) *Parser {
	if logger == nil {
		logger = slog.Default()
	}

	return &Parser{
		logger: logger,
	}
}

func (p *Parser) Parse(input string) (Query, error) {
	queryParts := strings.Fields(input)
	if len(queryParts) == 0 {
		p.logger.Debug("no query parts", slog.String("query", input))
		return Query{}, ErrInvalidQuery
	}

	command := queryParts[0]
	commandId := commandNameToId(command)
	if commandId == UnknownCommandId {
		p.logger.Debug("unknown command id", slog.String("query", input))
		return Query{}, ErrInvalidCommand
	}

	argumentsNumber := commandArgumentsNumber(commandId)
	query := NewQuery(commandId, queryParts[1:])
	if argumentsNumber != len(query.Arguments()) {
		p.logger.Debug("invalid arguments number", slog.String("query", input))
		return Query{}, ErrInvalidArgsNumber
	}

	return query, nil
}
