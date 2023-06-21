package logger

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	pb "github.com/amb1s1/gonetconfig/proto"
)

var (
	cisco = `
!
service timestamps log datetime msec show-timezone
logging trap informational
logging archive
 device harddisk
 severity informational
 file-size 10
 archive-size 100
 archive-length 52
!
logging console disable
logging monitor informational
logging buffered 10000000
logging buffered informational
logging facility local1
logging 192.168.1.1 vrf default severity debugging port default
logging 192.168.1.2 vrf default severity debugging port default
logging source-interface loopback0
logging hostnameprefix rt01.foo01
!
`
)

func TestLogger(t *testing.T) {
	tests := []struct {
		name       string
		device     *pb.Device
		wantConfig *pb.ConfigFeature
	}{
		{
			name: "Logger: not supported vendor => success.",
			device: &pb.Device{
				Name:   "rt01.foo01",
				Vendor: pb.Vendor_VD_ARISTA,
			},
			wantConfig: &pb.ConfigFeature{
				Name:          "logger",
				Version:       1,
				Configuration: "",
			},
		},
		{
			name: "Logger: Cisco config => success.",
			device: &pb.Device{
				Name:   "rt01.foo01",
				Vendor: pb.Vendor_VD_CISCO,
			},
			wantConfig: &pb.ConfigFeature{
				Name:          "logger",
				Version:       1,
				Configuration: cisco,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			in := &pb.ConfigGenRequest{
				Device: test.device,
			}
			out := &pb.ConfigGenResponse{ConfigFeature: []*pb.ConfigFeature{}}
			gen := New()
			gotConfig := gen.Render(ctx, in, out)
			if diff := cmp.Diff(test.wantConfig, gotConfig); diff != "" {
				t.Errorf("%v: unexpected diff (-want, +got):\n%v", test.name, diff)
			}
		})
	}
}
