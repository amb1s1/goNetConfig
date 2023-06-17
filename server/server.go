package server

import (
	"fmt"
	"net"

	"github.com/amb1s1/gonetconfig/service"
	"google.golang.org/grpc"

	pb "github.com/amb1s1/gonetconfig/proto"
)

const (
	protocol = "tcp"
	address  = "localhost:50051"
)

// Start starts a new server.
func Start() error {
	lis, err := net.Listen(protocol, address)
	if err != nil {
		return fmt.Errorf("Listen(), failed to open port: %w", err)
	}
	s := grpc.NewServer()
	// Register your gRPC service handlers here
	// pb.RegisterGoNetConfigServiceServer(s, NewServer())
	pb.RegisterGoNetConfigServiceServer(s, service.NewServer())
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}
