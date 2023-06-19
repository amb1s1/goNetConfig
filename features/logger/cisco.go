package logger

import (
	"bytes"
	"log"
	"text/template"

	pb "github.com/amb1s1/gonetconfig/proto"
	"github.com/golang/protobuf/proto"
)

var (
	version      int32 = 1
	mgtInterface       = "loopback0"
	loggerIPS          = []string{"192.168.1.1", "192.168.1.2"}
)

func ciscoRender(in *pb.ConfigGenRequest, out *pb.ConfigGenResponse) string {
	p := params{
		Hostname:            *in.Device.Name,
		ManagementInterface: mgtInterface,
		LoggerServerIPS:     loggerIPS,
	}
	switch in.Device.GetModel() {
	default:
		return defaultRender(out, p)
	}
}

func defaultRender(out *pb.ConfigGenResponse, p params) string {
	t, err := template.New(featureName).Parse(ciscoxrTemplate)
	if err != nil {
		log.Printf("New(), failed to parse template for %s feature", featureName)
	}
	var tpl bytes.Buffer
	err = t.Execute(&tpl, p)
	if err != nil {
		log.Printf("Execute(), failed to Execute template for %s feature", featureName)
	}
	if tpl.String() == "" {
		out.Status = proto.String("template not generated")
	}
	out.Status = proto.String("ok")

	return tpl.String()
}
