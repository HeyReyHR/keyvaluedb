package network

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"
	"time"
)

type TCPHandler = func(context.Context, []byte) []byte

type TCPServer struct {
	listener net.Listener

	idleTimeout    time.Duration
	bufferSize     int
	maxConnections int

	logger *slog.Logger
}

func NewTCPServer(address string, logger *slog.Logger, options ...TCPServerOption) (*TCPServer, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is invalid")
	}
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	server := &TCPServer{
		listener: listener,
		logger:   logger,
	}
	for _, option := range options {
		option(server)
	}

	if server.bufferSize == 0 {
		server.bufferSize = 4 << 20
	}

	return server, nil
}

func (s *TCPServer) HandleQueries(ctx context.Context, handler TCPHandler) {
	var wg sync.WaitGroup
	wg.Go(func() {
		for {
			connection, err := s.listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}
				s.logger.Error("failed to acc", slog.String("err", err.Error()))
				continue
			}
			go func(conn net.Conn) {
				s.handleConnection(ctx, connection, handler)
			}(connection)
		}
	})
	<-ctx.Done()
	s.listener.Close()

	wg.Wait()
}

func (s *TCPServer) handleConnection(ctx context.Context, connection net.Conn, handler TCPHandler) {
	request := make([]byte, s.bufferSize)

	for {
		if s.idleTimeout != 0 {
			if err := connection.SetReadDeadline(time.Now().Add(s.idleTimeout)); err != nil {
				s.logger.Warn("failed to set read deadline", slog.Any("error", err))
				break
			}
		}

		count, err := connection.Read(request)
		if err != nil && err != io.EOF {
			s.logger.Warn("failed to read data", slog.String("address", connection.RemoteAddr().String()), slog.Any("error", err))
			break
		} else if count == s.bufferSize {
			s.logger.Warn("small buffer size", slog.Int("buffer_size", s.bufferSize))
			break
		}

		if s.idleTimeout != 0 {
			if err := connection.SetWriteDeadline(time.Now().Add(s.idleTimeout)); err != nil {
				s.logger.Warn("failed to write deadline", slog.Any("error", err))
				break
			}
		}

		response := handler(ctx, request[:count])
		if _, err := connection.Write(response); err != nil {
			s.logger.Warn("failed to write data", slog.String("address", connection.RemoteAddr().String()), slog.Any("error", err))
			break
		}

	}
}
