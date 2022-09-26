package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"

	"github.com/amb1s1/gonetconfig/base"
	"github.com/amb1s1/gonetconfig/features/aaa"
	"google.golang.org/grpc"

	pb "github.com/amb1s1/gonetconfig/proto"
)

const (
	protocol = "tcp"
	address  = "localhost:50051"
)

var (
	registry     = newGeneratorsRegistry()
	allGenerator = []func() base.GeneratorInt{
		aaa.New,
	}
)

type server struct {
	pb.UnimplementedGoNetConfigServiceServer
}

type generatorsRegistry struct {
	generators    map[reflect.Type]base.GeneratorInt
	generatorType []reflect.Type
}

func newGeneratorsRegistry() *generatorsRegistry {
	return &generatorsRegistry{
		generators: make(map[reflect.Type]base.GeneratorInt),
	}
}
func main() {
	if err := registry.registerFeatures(); err != nil {
		log.Fatalln(err)
	}
	if err := star(); err != nil {
		log.Fatalln(err)
	}
}

func (s *server) GetConfigGen(ctx context.Context, in *pb.GetConfigGenRequest) (*pb.GetConfigGenResponse, error) {
	response := &pb.GetConfigGenResponse{}
	for _, g := range registry.generatorType {
		response = &pb.GetConfigGenResponse{
			ConfigFeature: registry.generators[g].Render(ctx),
		}
	}
	return response, nil
}

func star() error {
	lis, err := net.Listen(protocol, address)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to open port, error: %v", err))
	}
	s := grpc.NewServer()
	pb.RegisterGoNetConfigServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		return errors.New(fmt.Sprintf("failed to serve: %v", err))
	}
	return nil
}

func (g *generatorsRegistry) registerFeatures() error {
	for _, generator := range allGenerator {
		kind := reflect.TypeOf(generator())
		if _, exists := g.generators[kind]; exists {
			return fmt.Errorf("generator already exist, generator %v", kind)
		}
		g.generators[kind] = generator()
		g.generatorType = append(g.generatorType, kind)
	}
	return nil
}
