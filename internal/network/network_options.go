package network

import "time"

type TCPClientOption func(*TCPClient)

type TCPServerOption func(*TCPServer)

const defaultBufferSize = 10

func WithClientBufferSize(buffer int) TCPClientOption {
	return func(c *TCPClient) {
		c.bufferSize = buffer
	}
}

func WithClientIdleTimeout(timeout time.Duration) TCPClientOption {
	return func(c *TCPClient) {
		c.idleTimeout = timeout
	}
}

func WithServerBufferSize(buffer uint) TCPServerOption {
	return func(s *TCPServer) {
		s.bufferSize = int(buffer)
	}
}

func WithServerMaxConnections(maxConns uint) TCPServerOption {
	return func(s *TCPServer) {
		s.maxConnections = int(maxConns)
	}
}

func WithServerIdleTimeout(timeout time.Duration) TCPServerOption {
	return func(s *TCPServer) {
		s.idleTimeout = timeout
	}
}
