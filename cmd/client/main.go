package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/HeyReyHR/keyvaluedb/internal/common"
	"github.com/HeyReyHR/keyvaluedb/internal/network"
)

func main() {
	address := flag.String("address", "localhost:3223", "Address of db")
	idleTimeout := flag.Duration("idle_timeout", time.Minute, "Idle timeout for connection")
	maxMessageSizeStr := flag.String("max_message_size", "4KB", "Max message size for connection")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	maxMessageSize, err := common.ParseSize(*maxMessageSizeStr)
	if err != nil {
		logger.Error("failed to parse max message size", slog.Any("error", err))
		os.Exit(1)
	}

	var options []network.TCPClientOption
	options = append(options, network.WithClientIdleTimeout(*idleTimeout))
	options = append(options, network.WithClientBufferSize(maxMessageSize))

	reader := bufio.NewReader(os.Stdin)
	client, err := network.NewTCPClient(*address, options...)
	if err != nil {
		logger.Error("failed to connect to server", slog.Any("error", err))
		os.Exit(1)
	}

	for {
		fmt.Print("[db] ")
		request, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				logger.Info("stdin closed, exiting")
				os.Exit(0)
			}
			logger.Error("failed to read query", slog.Any("error", err))
			continue
		}
		response, err := client.Send([]byte(request))
		if err != nil {
			if isConnectionDead(err) {
				logger.Error("connection was closed", slog.Any("error", err))
				os.Exit(1)
			}
			logger.Error("failed to send query", slog.Any("error", err))
			continue
		}

		fmt.Println(string(response))
	}
}

func isConnectionDead(err error) bool {
	if errors.Is(err, io.EOF) {
		return true
	}
	if errors.Is(err, syscall.EPIPE) {
		return true
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}

	return false
}
