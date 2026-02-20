package service

import (
	"context"
	"fmt"
	"reflect"

	"github.com/amb1s1/gonetconfig/base"
	"github.com/amb1s1/gonetconfig/entity"
	"github.com/amb1s1/gonetconfig/features/aaa"
	"github.com/amb1s1/gonetconfig/features/logger"

	pb "github.com/amb1s1/gonetconfig/proto"
)

var (
	// Registry is a map of generators, keyed by their type.
	Registry = newGeneratorsRegistry()
	// allGenerator is a list of all available generators.
	allGenerator = []func() base.GeneratorInt{
		aaa.New,
		logger.New,
	}
)

type generatorsRegistry struct {
	generators     map[reflect.Type]base.GeneratorInt
	entityResolved *entity.ResolvedEntity
}

type server struct {
	pb.UnimplementedGoNetConfigServiceServer
}

// NewServer creates a new server.
func NewServer() *server {
	return &server{}
}

func newGeneratorsRegistry() *generatorsRegistry {
	return &generatorsRegistry{
		generators: make(map[reflect.Type]base.GeneratorInt),
	}
}

// SetEntity configures the registry with a resolved entity, applying feature
// filtering and parameter injection to all registered generators.
func (g *generatorsRegistry) SetEntity(resolved *entity.ResolvedEntity) {
	g.entityResolved = resolved
	// Always reset params first to avoid leaking values across entity changes.
	for _, gen := range g.generators {
		gen.SetParams(nil)
	}

	if resolved == nil {
		return
	}

	for _, gen := range g.generators {
		featureName := gen.Generators().Name
		if params, ok := resolved.Params[featureName]; ok {
			gen.SetParams(params)
		}
	}
}

// isFeatureEnabled checks if a feature is enabled for the current entity.
// If no entity is configured, all features are enabled.
func (g *generatorsRegistry) isFeatureEnabled(featureName string) bool {
	if g.entityResolved == nil {
		return true
	}
	for _, f := range g.entityResolved.Features {
		if f == featureName {
			return true
		}
	}
	return false
}

// GetConfigGen gets a config generator for the specified device.
func (s *server) GetConfigGen(ctx context.Context, in *pb.ConfigGenRequest) (*pb.ConfigGenResponse, error) {
	out := &pb.ConfigGenResponse{
		ConfigFeature: []*pb.ConfigFeature{},
	}
	for _, g := range Registry.generators {
		if !Registry.isFeatureEnabled(g.Generators().Name) {
			continue
		}
		if generatorSupported(g.Generators(), in) {
			out.ConfigFeature = append(out.ConfigFeature, g.Render(ctx, in, out))
		}
	}
	return out, nil
}

// RegisterFeatures registers all available generators with the registry.
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
