package service

import (
	"context"
	"fmt"
	"reflect"

	"github.com/amb1s1/gonetconfig/base"
	"github.com/amb1s1/gonetconfig/features/aaa"
	"github.com/amb1s1/gonetconfig/features/logger"

	pb "github.com/amb1s1/gonetconfig/proto"
)

var (
	Registry     = newGeneratorsRegistry()
	allGenerator = []func() base.GeneratorInt{
		aaa.New,
		logger.New,
	}
)

type generatorsRegistry struct {
	generators map[reflect.Type]base.GeneratorInt
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
	out := &pb.ConfigGenResponse{
		ConfigFeature: []*pb.ConfigFeature{},
	}
	for _, g := range Registry.generators {
		if generatorSupported(g.Generators(), in) {
			out.ConfigFeature = append(out.ConfigFeature, g.Render(ctx, in, out))
		}
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
	}
	return nil
}

func generatorSupported(gen *base.Generator, in *pb.ConfigGenRequest) bool {
	if !isVendorSupported(gen, &in.Device.Vendor) {
		return false
	}
	if !isModelSupported(gen, &in.Device.Model) {
		return false
	}
	return true
}
func isVendorSupported(gen *base.Generator, vendor *pb.Vendor) bool {
	return vendor == nil || gen.SupportedVendor[*vendor] || *vendor == pb.Vendor_VD_UNKNOW
}

func isModelSupported(gen *base.Generator, model *pb.Model) bool {
	return model == nil || gen.SupportedModel[*model] || *model == pb.Model_MD_UNKNOW
}
