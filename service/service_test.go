package service

import (
	"context"
	"testing"

	pb "github.com/amb1s1/gonetconfig/proto"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
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
	ciscoLogger = `
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

func TestService(t *testing.T) {
	tests := []struct {
		name      string
		device    *pb.Device
		wantRespo *pb.ConfigGenResponse
	}{
		{
			name: "service: config cisco",
			device: &pb.Device{
				Name:   "rt01.foo01",
				Vendor: pb.Vendor_VD_CISCO,
			},
			wantRespo: &pb.ConfigGenResponse{
				Status: "ok",
				ConfigFeature: []*pb.ConfigFeature{
					{
						Name:          "aaa",
						Version:       1,
						Configuration: ciscoAAA,
					},
					{
						Name:          "logger",
						Version:       1,
						Configuration: ciscoLogger,
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.TODO()
			request := &pb.ConfigGenRequest{
				Device: test.device,
			}
			registry := Registry

			// Register the available generators
			err := registry.RegisterFeatures()
			assert.NoError(t, err)

			// Create the server
			server := NewServer()

			// Call the GetConfigGen function
			response, err := server.GetConfigGen(ctx, request)

			if err != nil {
				t.Errorf("GetConfigGen() unexpected error: %v", err)
			}
			if response == nil {
				t.Errorf("GetConfigGen() response is unexpected nil")
			}
			if want, got := test.wantRespo.Status, response.Status; want != got {
				t.Errorf("unexpected status, want: %s, got: %s", want, got)
			}
			if diff := cmp.Diff(test.wantRespo.ConfigFeature, response.ConfigFeature); diff != "" {
				t.Errorf("unexpected diff(want- got+): %v", diff)
			}
		})
	}
}
