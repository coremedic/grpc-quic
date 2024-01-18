# gRPC over QUIC

The Go language implementation of gRPC over QUIC.

Largely based on 'https://github.com/sssgun/grpc-quic' updated for quic-go

## Installation

```shell
go get github.com/coremedic/grpc-quic
```

## Usage

### Import module
```go
grpcquic "github.com/coremedic/grpc-quic"
```

### Server

```go
func main() {
	// Generate TLS config for server
	log.Println("Generating TLS config...")
	tlsConfig, err := generateTLSConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Start gRPC QUIC listener and service
	log.Println("Starting QUIC gRPC Server...")
	quicListener, err := quic.ListenAddr("127.0.0.1:1848", tlsConfig, nil)
	if err != nil {
		log.Fatal(err)
	}

	grpcQuicListener := grpcquic.Listen(*quicListener)

	grpcServer := grpc.NewServer()
	pb.RegisterExampleServiceServer(grpcServer, &grpcServiceServer{})
	log.Printf("gRPC-QUIC: listening at %v\n", grpcQuicListener.Addr())

	// Accept incoming gRPC requests
	if err := grpcServer.Serve(grpcQuicListener); err != nil {
		log.Fatal(err)
	}
}
```

### Client

```go
func main() {
	// Create TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"grpc-quic-example"},
	}
	creds := grpcquic.NewCredentials(tlsConfig)

	// Connect to gRPC Service Server
	dialer := grpcquic.NewQuicDialer(tlsConfig)
	grpcOpts := []grpc.DialOption{
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(creds),
	}
	conn, err := grpc.Dial("127.0.0.1:1848", grpcOpts...)
	if err != nil {
		log.Fatal(err)
	}

	// Close connection at end of function
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

	// Create gRPC client
	grpcClient := pb.NewExampleServiceClient(conn)

	// Send gRPC request
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req, err := grpcClient.Example(ctx, &pb.ExampleRequest{Msg: "Ayooooo"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[gRPC]: %v\n", req.GetMsg())
}
```