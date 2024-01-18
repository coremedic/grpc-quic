package grpcquic

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	"google.golang.org/grpc/credentials"
)

// ErrInvalidConnType is returned when a connection is not of type QuicConn.
var ErrInvalidConnType = fmt.Errorf("connection must be of type QuicConn")

// authInfo facade structure for gRPC Credentials using QUIC
type authInfo struct {
	conn *QuicConn
}

// Implementation of AuthType
// https://pkg.go.dev/google.golang.org/grpc@v1.42.0/credentials#pkg-types
func (ai *authInfo) AuthType() string {
	return "quic-tls"
}

// newAuthInfo creates a new AuthInfo object from a QuicConn object
func newAuthInfo(conn *QuicConn) credentials.AuthInfo {
	return &authInfo{conn: conn}
}

/*
Credentials for gRPC over QUIC
https://pkg.go.dev/google.golang.org/grpc@v1.42.0/credentials#TransportCredentials
*/
type Credentials struct {
	tlsConfig  *tls.Config
	quicConn   bool
	serverName string
	grpcCreds  credentials.TransportCredentials
}

/*
QUIC gRPC implementation of ClientHandshake

ClientHandshake handles the client-side authentication handshake for the
specified authentication protocol. ClientHandshake returns an authenicated connection
and relevant authentication information.

Read more here: https://pkg.go.dev/google.golang.org/grpc@v1.42.0/credentials#TransportCredentials
*/
func (creds *Credentials) ClientHandshake(ctx context.Context, authority string, conn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	if c, ok := conn.(*QuicConn); ok {
		creds.quicConn = true
		return conn, newAuthInfo(c), nil
	}
	return nil, nil, ErrInvalidConnType
}

/*
QUIC gRPC implementation of ServerHandshake

ServerHandshake handles the server-side authentication handshake for the
specified authentication protocol. ServerHandshake returns an authenicated connection
and relevant authentication information.

Read more here: https://pkg.go.dev/google.golang.org/grpc@v1.42.0/credentials#TransportCredentials
*/
func (creds *Credentials) ServerHandshake(conn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	if c, ok := conn.(*QuicConn); ok {
		creds.quicConn = true
		return conn, newAuthInfo(c), nil
	}
	return nil, nil, ErrInvalidConnType
}

// Info provides ProtocolInfo of Credentials
func (creds *Credentials) Info() credentials.ProtocolInfo {
	// We only handle QUIC connections
	if !creds.quicConn {
		return creds.grpcCreds.Info()
	}

	return credentials.ProtocolInfo{
		ProtocolVersion:  "grpc-go-quic/1.0.0",
		SecurityProtocol: "quic-tls",
		ServerName:       creds.serverName,
	}
}

/*
QUIC gRPC facade for OverrideServerName
*/
// Deprecated: use grpc.WithAuthority instead.
func (creds *Credentials) OverrideServerName(name string) error {
	creds.serverName = name
	return creds.grpcCreds.OverrideServerName(name)
}

// Clone makes a copy of Credentials
func (creds *Credentials) Clone() credentials.TransportCredentials {
	return &Credentials{
		tlsConfig:  creds.tlsConfig.Clone(),
		quicConn:   creds.quicConn,
		serverName: creds.serverName,
		grpcCreds:  creds.grpcCreds.Clone(),
	}
}

// Create new credentials object from tls config
func NewCredentials(tlsConfig *tls.Config) credentials.TransportCredentials {
	grpcCreds := credentials.NewTLS(tlsConfig)
	return &Credentials{
		grpcCreds: grpcCreds,
		tlsConfig: tlsConfig,
	}
}
