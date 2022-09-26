package service

import (
	"context"
	"errors"
	"fmt"
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
	Registry     = newGeneratorsRegistry()
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

func (s *server) GetConfigGen(ctx context.Context, in *pb.GetConfigGenRequest) (*pb.GetConfigGenResponse, error) {
	response := &pb.GetConfigGenResponse{}
	for _, g := range Registry.generatorType {
		response = &pb.GetConfigGenResponse{
			ConfigFeature: Registry.generators[g].Render(ctx),
		}
	}
	return response, nil
}

func (g *generatorsRegistry) RegisterFeatures() error {
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

func Start() error {
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
