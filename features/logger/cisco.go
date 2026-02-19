package logger

import (
	"bytes"
	"log"
	"text/template"

	pb "github.com/amb1s1/gonetconfig/proto"
)

var (
	version            int32 = 1
	defaultMgtInterface      = "loopback0"
	defaultLoggerIPs         = []string{"192.168.1.1", "192.168.1.2"}
)

func (g *Generator) ciscoRender(in *pb.ConfigGenRequest, out *pb.ConfigGenResponse) string {
	p := params{
		Hostname:            in.Device.Name,
		ManagementInterface: g.getManagementInterface(),
		LoggerServerIPS:     g.getLoggerServerIPs(),
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
		out.Status = "template not generated"
	}
	out.Status = "ok"
	return tpl.String()
}
