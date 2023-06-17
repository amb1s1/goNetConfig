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

type generatorsRegistry struct {
	generators    map[reflect.Type]base.GeneratorInt
	generatorType []reflect.Type
}

type server struct {
	pb.UnimplementedGoNetConfigServiceServer
}

func NewServer() *server {
	return &server{}
}

func newGeneratorsRegistry() *generatorsRegistry {
	return &generatorsRegistry{
		generators: make(map[reflect.Type]base.GeneratorInt),
	}
}

func (s *server) GetConfigGen(ctx context.Context, in *pb.ConfigGenRequest) (*pb.ConfigGenResponse, error) {
	out := &pb.ConfigGenResponse{}
	for _, g := range Registry.generatorType {
		out.ConfigFeature = Registry.generators[g].Render(ctx, in, out)
	}
	return out, nil
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
