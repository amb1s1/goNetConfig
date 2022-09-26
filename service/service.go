package service

import (
	"context"
	"fmt"
	"reflect"

	"github.com/amb1s1/gonetconfig/base"
	"github.com/amb1s1/gonetconfig/features/aaa"

	pb "github.com/amb1s1/gonetconfig/proto"
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

func NewServer() *server {
	return &server{}
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
