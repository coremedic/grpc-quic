package main

import (
	"context"
	"crypto/tls"
	"log"
	"time"

	grpcquic "github.com/coremedic/grpc-quic"
	pb "github.com/coremedic/grpc-quic/example/protobuf"
	"google.golang.org/grpc"
)

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
