package base

import (
	"context"

	pb "github.com/amb1s1/gonetconfig/proto"
)

type Generator struct {
	Name   string
	Vendor string
	Model  string
}

type GeneratorInt interface {
	Render(context.Context) *pb.ConfigFeature
	Generators() *Generator
}
