package grpcquic

import (
	"context"
	"net"
	"time"

	"github.com/quic-go/quic-go"
)

// QuicConn implements net.Conn using a QUIC connection and stream
type QuicConn struct {
	connection *quic.Conn
	stream     *quic.Stream
}

// Read reads data from the QUIC stream. Returns the number of bytes read and any error encountered
func (qc *QuicConn) Read(b []byte) (int, error) {
	if qc.stream == nil {
		return 0, net.ErrClosed
	}
	return qc.stream.Read(b)
}

// Write writes data to the QUIC stream. Returns the number of bytes written and any error encountered
func (qc *QuicConn) Write(b []byte) (int, error) {
	if qc.stream == nil {
		return 0, net.ErrClosed
	}
	return qc.stream.Write(b)
}

// Close closes both the QUIC stream and connection
func (qc *QuicConn) Close() error {
	if qc.stream != nil {
		if err := qc.stream.Close(); err != nil {
			return err
		}
	}
	return qc.connection.CloseWithError(0, "")
}

// LocalAddr returns the local network address of the QUIC connection
func (qc *QuicConn) LocalAddr() net.Addr {
	return qc.connection.LocalAddr()
}

// RemoteAddr returns the remote network address of the QUIC connection
func (qc *QuicConn) RemoteAddr() net.Addr {
	return qc.connection.RemoteAddr()
}

// SetDeadline sets the read and write deadlines for the QUIC stream
func (qc *QuicConn) SetDeadline(t time.Time) error {
	return qc.stream.SetDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls on the QUIC stream
func (qc *QuicConn) SetReadDeadline(t time.Time) error {
	return qc.stream.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls on the QUIC stream
func (qc *QuicConn) SetWriteDeadline(t time.Time) error {
	return qc.stream.SetWriteDeadline(t)
}

// NewQuicConn creates a new QuicConn with an open QUIC stream
func NewQuicConn(connection *quic.Conn) (net.Conn, error) {
	stream, err := connection.OpenStreamSync(context.Background())
	if err != nil {
		return nil, err
	}
	return &QuicConn{connection: connection, stream: stream}, nil
}
