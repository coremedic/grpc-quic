package grpcquic

import (
	"context"
	"net"

	"github.com/quic-go/quic-go"
)

type QuicListener struct {
	ql quic.Listener
}

// Accept waits for and returns the next connection to the listener
func (ql *QuicListener) Accept() (net.Conn, error) {
	conn, err := ql.ql.Accept(context.Background())
	if err != nil {
		return nil, err
	}

	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		return nil, err
	}

	return &QuicConn{conn, stream}, nil
}

// Close closes the listener
func (ql *QuicListener) Close() error {
	return ql.ql.Close()
}

// Addr returns the listeners network address
func (ql *QuicListener) Addr() net.Addr {
	return ql.ql.Addr()
}

func Listen(ql quic.Listener) net.Listener {
	return &QuicListener{ql}
}
