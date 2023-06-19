package logger

import (
	"context"

	"github.com/amb1s1/gonetconfig/base"
	"github.com/golang/protobuf/proto"

	pb "github.com/amb1s1/gonetconfig/proto"
)

type Generator struct {
	baseGen *base.Generator
}

const (
	featureName = "logger"
)

var (
	generator = &base.Generator{
		Name:            featureName,
		SupportedVendor: map[pb.Vendor]bool{pb.Vendor_VD_CISCO: true},
	}
)

type params struct {
	Hostname            string
	ManagementInterface string
	LoggerServerIPS     []string
}

func New() base.GeneratorInt {
	gen := &Generator{}
	gen.baseGen = generator
	return gen
}

func (g *Generator) Render(ctx context.Context, in *pb.ConfigGenRequest, out *pb.ConfigGenResponse) *pb.ConfigFeature {
	rendering := renderTemplate(in, out)
	return &pb.ConfigFeature{
		Name:          proto.String(featureName),
		Version:       proto.Int32(1),
		Configuration: proto.String(rendering),
	}
}

func (g *Generator) Generators() *base.Generator {
	return g.baseGen
}

func renderTemplate(in *pb.ConfigGenRequest, out *pb.ConfigGenResponse) string {
	var render string
	switch in.GetDevice().GetVendor() {
	case *pb.Vendor_VD_CISCO.Enum():
		render = ciscoRender(in, out)
	}
	return render
}

func isSupported(model string) bool {
	return true
}
