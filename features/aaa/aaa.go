package aaa

import (
	"bytes"
	"context"
	"log"
	"text/template"

	"github.com/amb1s1/gonetconfig/base"
	"github.com/golang/protobuf/proto"

	pb "github.com/amb1s1/gonetconfig/proto"
)

type Generator struct {
	baseGen *base.Generator
}

var (
	aaaName   = "aaa"
	generator = &base.Generator{
		Name:   "aaa",
		Vendor: "ciscoxr",
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

func (g *Generator) Render(ctx context.Context) *pb.ConfigFeature {
	rendering := renderTemplate()
	return &pb.ConfigFeature{
		Name:          proto.String("aaa"),
		Version:       proto.Int32(1),
		Configuration: proto.String(rendering),
	}
}

func (g *Generator) Generators() *base.Generator {
	return g.baseGen
}

func renderTemplate() string {
	t, err := template.New("aaa").Parse(ciscoxrTemplate)
	if err != nil {
		log.Panic(err)
	}
	var tpl bytes.Buffer
	err = t.Execute(&tpl, params{Secret: "password123"})
	if err != nil {
		log.Panic(err)
	}
	return tpl.String()
}

func isSupported(model string) bool {
	return true
}
