package grpcquic

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/quic-go/quic-go"
)

// Default QUIC config
var defaultQuicConfig = quic.Config{
	KeepAlivePeriod: 10 * time.Second,
}

// NewQuicDialer returns a new Quic Dialer function
func NewQuicDialer(tlsConfig *tls.Config) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, addr string) (net.Conn, error) {
		conn, err := quic.DialAddr(ctx, addr, tlsConfig, &defaultQuicConfig)
		if err != nil {
			return nil, err
		}
		return NewQuicConn(conn)
	}
}
