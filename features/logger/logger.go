package logger

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

// getManagementInterface returns the entity param or falls back to the default.
func (g *Generator) getManagementInterface() string {
	if g.entityParams != nil {
		if v, ok := g.entityParams["management_interface"]; ok {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}
	return defaultMgtInterface
}

// getLoggerServerIPs returns the entity param or falls back to the default.
func (g *Generator) getLoggerServerIPs() []string {
	if g.entityParams != nil {
		if v, ok := g.entityParams["logger_server_ips"]; ok {
			if slice, ok := v.([]interface{}); ok {
				ips := make([]string, 0, len(slice))
				for _, item := range slice {
					if s, ok := item.(string); ok {
						ips = append(ips, s)
					}
				}
				return ips
			}
		}
	}
	return defaultLoggerIPs
}

func (g *Generator) renderTemplate(in *pb.ConfigGenRequest, out *pb.ConfigGenResponse) string {
	var render string
	switch in.GetDevice().GetVendor() {
	case pb.Vendor_VD_CISCO:
		render = g.ciscoRender(in, out)
	}
	return render
}

func isSupported(model string) bool {
	return true
}
