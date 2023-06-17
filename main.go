package main

import (
	"fmt"
	"log"
	"net"

	"github.com/amb1s1/gonetconfig/service"
	"google.golang.org/grpc"

	pb "github.com/amb1s1/gonetconfig/proto"
)

const (
	protocol = "tcp"
	address  = "localhost:50051"
)

func main() {
	if err := service.Registry.RegisterFeatures(); err != nil {
		log.Fatalf("RegisterFeatures(), failed to register features: %v", err)
	}
	if err := start(); err != nil {
		log.Fatalf("start(), failed to start the server: %v", err)
	}
}

func start() error {
	lis, err := net.Listen(protocol, address)
	if err != nil {
		return fmt.Errorf("Listen(), failed to open port: %w", err)
	}
	s := grpc.NewServer()
	pb.RegisterGoNetConfigServiceServer(s, service.NewServer())
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}
