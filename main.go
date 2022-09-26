package main

import (
	"errors"
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
		log.Fatalln(err)
	}
	if err := start(); err != nil {
		log.Fatalln(err)
	}
}

func start() error {
	lis, err := net.Listen(protocol, address)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to open port, error: %v", err))
	}
	s := grpc.NewServer()
	pb.RegisterGoNetConfigServiceServer(s, service.NewServer())
	if err := s.Serve(lis); err != nil {
		return errors.New(fmt.Sprintf("failed to serve: %v", err))
	}
	return nil
}
