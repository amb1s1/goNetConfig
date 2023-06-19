package aaa

import (
	"bytes"
	"log"
	"text/template"

	pb "github.com/amb1s1/gonetconfig/proto"
)

var (
	secret        = "password!123"
	version int32 = 1
)

func ciscoRender(in *pb.ConfigGenRequest, out *pb.ConfigGenResponse) string {
	switch in.Device.GetModel() {
	default:
		return defaultRender(out)
	}
}

func defaultRender(out *pb.ConfigGenResponse) string {
	t, err := template.New(featureName).Parse(ciscoxrTemplate)
	if err != nil {
		log.Printf("New(), failed to parse template for %s feature", featureName)
	}
	var tpl bytes.Buffer
	err = t.Execute(&tpl, params{Secret: secret})
	if err != nil {
		log.Printf("Execute(), failed to Execute template for %s feature", featureName)
	}
	if tpl.String() == "" {
		out.Status = "template not generated"
	}
	out.Status = "ok"

	return tpl.String()
}
