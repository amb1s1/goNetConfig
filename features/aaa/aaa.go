package aaa

import (
	"context"

	"github.com/amb1s1/gonetconfig/base"

	pb "github.com/amb1s1/gonetconfig/proto"
)

type Generator struct {
	baseGen *base.Generator
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
	rendering := renderTemplate(in, out)
	return &pb.ConfigFeature{
		Name:          featureName,
		Version:       1,
		Configuration: rendering,
	}
}

func (g *Generator) Generators() *base.Generator {
	return g.baseGen
}

func renderTemplate(in *pb.ConfigGenRequest, out *pb.ConfigGenResponse) string {
	var render string
	switch in.GetDevice().GetVendor() {
	case pb.Vendor_VD_CISCO:
		render = ciscoRender(in, out)
	}
	return render
}

func isSupported(model string) bool {
	return true
}
