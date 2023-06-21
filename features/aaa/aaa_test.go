package aaa

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	pb "github.com/amb1s1/gonetconfig/proto"
)

var (
	ciscoAAA = `
!
usergroup priv15
 taskgroup root-lr
 taskgroup cisco-support
!
username netops
 group root-lr
 group cisco-support
 secret 5 password!123
!
`
)

func TestAAA(t *testing.T) {
	tests := []struct {
		name       string
		device     string
		vendor     pb.Vendor
		model      pb.Model
		wantConfig *pb.ConfigFeature
		wantStatus
	}{
		{
			name:   "AAA: not supported vendor => success.",
			device: "rt01.foo01",
			vendor: pb.Vendor_VD_ARISTA,
			wantConfig: &pb.ConfigFeature{
				Name:          "aaa",
				Version:       1,
				Configuration: "",
			},
		},
		{
			name:   "AAA: supported vendor, but not supported model => success.",
			device: "rt01.foo01",
			vendor: pb.Vendor_VD_ARISTA,
			model:  pb.Model_MD_C3560CX,
			wantConfig: &pb.ConfigFeature{
				Name:          "aaa",
				Version:       1,
				Configuration: "",
			},
		},
		{
			name:   "AAA: Cisco triple aaa config => success.",
			device: "rt01.foo01",
			vendor: pb.Vendor_VD_CISCO,
			wantConfig: &pb.ConfigFeature{
				Name:          "aaa",
				Version:       1,
				Configuration: ciscoAAA,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			in := &pb.ConfigGenRequest{
				Device: &pb.Device{
					Name:   test.device,
					Vendor: test.vendor,
					Model:  test.model,
				},
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
