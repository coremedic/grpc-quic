package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"log"
	"math/big"

	grpcquic "github.com/coremedic/grpc-quic"
	pb "github.com/coremedic/grpc-quic/example/protobuf"
	"github.com/quic-go/quic-go"
	"google.golang.org/grpc"
)

type grpcServiceServer struct {
	pb.UnimplementedExampleServiceServer
}

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

// Implementation of Example rpc
func (gss *grpcServiceServer) Example(ctx context.Context, req *pb.ExampleRequest) (*pb.ExampleResponse, error) {
	log.Printf("[gRPC]: %v\n", req.GetMsg())
	return &pb.ExampleResponse{Msg: "what up"}, nil
}

// generateTLSConfig generates a new TLS configuration with an in-memory certificate and key pair.
func generateTLSConfig() (*tls.Config, error) {
	// Generate a new RSA key
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Printf("failed to generate RSA key: %s", err)
		return nil, err
	}

	// Create a certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		// Set other fields of the certificate as required
	}

	// Create a certificate using the template
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		log.Printf("failed to create certificate: %s", err)
		return nil, err
	}

	// Encode the certificate and key to PEM format
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	// Load the X509 key pair from PEM blocks
	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		log.Printf("failed to load X509 key pair from PEM: %s", err)
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"grpc-quic-example"},
	}, nil
}
