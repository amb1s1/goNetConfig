package aaa

import (
	"context"

	"github.com/amb1s1/gonetconfig/base"

	pb "github.com/amb1s1/gonetconfig/proto"
)

type Generator struct {
	baseGen      *base.Generator
	entityParams map[string]interface{}
}

const (
	featureName = "aaa"
)

var (
	generator = &base.Generator{
		Name:            featureName,
		SupportedVendor: map[pb.Vendor]bool{pb.Vendor_VD_CISCO: true},
	}
)

type params struct {
	Secret string
}

func New() base.GeneratorInt {
	gen := &Generator{}
	gen.baseGen = generator
	return gen
}

func (g *Generator) Render(ctx context.Context, in *pb.ConfigGenRequest, out *pb.ConfigGenResponse) *pb.ConfigFeature {
	rendering := g.renderTemplate(in, out)
	return &pb.ConfigFeature{
		Name:          featureName,
		Version:       1,
		Configuration: rendering,
	}
}

func (g *Generator) Generators() *base.Generator {
	return g.baseGen
}

// SetParams injects entity-specific parameters into the generator.
func (g *Generator) SetParams(p map[string]interface{}) {
	g.entityParams = p
}

// getSecret returns the entity param secret or falls back to the default.
func (g *Generator) getSecret() string {
	if g.entityParams != nil {
		if v, ok := g.entityParams["secret"]; ok {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}
	return defaultSecret
}

func (g *Generator) renderTemplate(in *pb.ConfigGenRequest, out *pb.ConfigGenResponse) string {
	var render string
	switch in.GetDevice().GetVendor() {
	case pb.Vendor_VD_CISCO:
		render = g.ciscoRender(in, out)
	}
	return render
}
