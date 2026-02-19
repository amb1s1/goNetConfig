package service

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/amb1s1/gonetconfig/entity"
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

func TestServiceWithEntity(t *testing.T) {
	// Save and restore global Registry
	origRegistry := Registry
	defer func() { Registry = origRegistry }()

	t.Run("entity filters features to aaa only", func(t *testing.T) {
		Registry = newGeneratorsRegistry()
		err := Registry.RegisterFeatures()
		assert.NoError(t, err)

		resolved := &entity.ResolvedEntity{
			Name:     "testEntity",
			Features: []string{"aaa"},
			Params: map[string]map[string]interface{}{
				"aaa": {"secret": "entity-secret-789"},
			},
		}
		Registry.SetEntity(resolved)

		ctx := context.TODO()
		request := &pb.ConfigGenRequest{
			Device: &pb.Device{
				Name:   "rt01.foo01",
				Vendor: pb.Vendor_VD_CISCO,
			},
		}

		server := NewServer()
		response, err := server.GetConfigGen(ctx, request)
		assert.NoError(t, err)

		// Should only have aaa, not logger
		assert.Len(t, response.ConfigFeature, 1)
		assert.Equal(t, "aaa", response.ConfigFeature[0].Name)
		// Should use entity secret
		assert.Contains(t, response.ConfigFeature[0].Configuration, "entity-secret-789")
	})

	t.Run("entity with all features and custom params", func(t *testing.T) {
		Registry = newGeneratorsRegistry()
		err := Registry.RegisterFeatures()
		assert.NoError(t, err)

		resolved := &entity.ResolvedEntity{
			Name:     "fullEntity",
			Features: []string{"aaa", "logger"},
			Params: map[string]map[string]interface{}{
				"aaa": {"secret": "full-secret"},
				"logger": {
					"management_interface": "loopback99",
					"logger_server_ips":    []interface{}{"1.2.3.4"},
				},
			},
		}
		Registry.SetEntity(resolved)

		ctx := context.TODO()
		request := &pb.ConfigGenRequest{
			Device: &pb.Device{
				Name:   "rt01.foo01",
				Vendor: pb.Vendor_VD_CISCO,
			},
		}

		server := NewServer()
		response, err := server.GetConfigGen(ctx, request)
		assert.NoError(t, err)

		assert.Len(t, response.ConfigFeature, 2)

		// Find features by name since map iteration order is not guaranteed
		featureMap := make(map[string]*pb.ConfigFeature)
		for _, f := range response.ConfigFeature {
			featureMap[f.Name] = f
		}

		aaaFeature := featureMap["aaa"]
		assert.NotNil(t, aaaFeature)
		assert.Contains(t, aaaFeature.Configuration, "full-secret")

		loggerFeature := featureMap["logger"]
		assert.NotNil(t, loggerFeature)
		assert.Contains(t, loggerFeature.Configuration, "loopback99")
		assert.Contains(t, loggerFeature.Configuration, "1.2.3.4")
		assert.NotContains(t, loggerFeature.Configuration, "192.168.1.1")
	})

	t.Run("nil entity uses all features with defaults", func(t *testing.T) {
		Registry = newGeneratorsRegistry()
		err := Registry.RegisterFeatures()
		assert.NoError(t, err)

		Registry.SetEntity(nil)

		ctx := context.TODO()
		request := &pb.ConfigGenRequest{
			Device: &pb.Device{
				Name:   "rt01.foo01",
				Vendor: pb.Vendor_VD_CISCO,
			},
		}

		server := NewServer()
		response, err := server.GetConfigGen(ctx, request)
		assert.NoError(t, err)

		// Should have both features with default values
		assert.Len(t, response.ConfigFeature, 2)
	})
}

func TestServiceWithSampleEntityConfigs(t *testing.T) {
	// Save and restore global Registry
	origRegistry := Registry
	defer func() { Registry = origRegistry }()

	tests := []struct {
		name            string
		sampleConfig    string
		wantFeatureCnt  int
		wantHasLogger   bool
		wantAAAContains string
		wantLogContains string
	}{
		{
			name:            "sample companyA enables aaa and logger",
			sampleConfig:    "entity_company_a.yml",
			wantFeatureCnt:  2,
			wantHasLogger:   true,
			wantAAAContains: "sample-a-secret",
			wantLogContains: "10.10.10.1",
		},
		{
			name:            "sample companyB enables only aaa",
			sampleConfig:    "entity_company_b.yml",
			wantFeatureCnt:  1,
			wantHasLogger:   false,
			wantAAAContains: "sample-b-secret",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			Registry = newGeneratorsRegistry()
			err := Registry.RegisterFeatures()
			assert.NoError(t, err)

			cfg, err := entity.Load(sampleConfigPath(t, tc.sampleConfig))
			assert.NoError(t, err)

			resolved, err := cfg.GetSelectedEntity()
			assert.NoError(t, err)

			Registry.SetEntity(resolved)

			request := &pb.ConfigGenRequest{
				Device: &pb.Device{
					Name:   "rt01.foo01",
					Vendor: pb.Vendor_VD_CISCO,
				},
			}

			server := NewServer()
			response, err := server.GetConfigGen(context.Background(), request)
			assert.NoError(t, err)
			assert.Len(t, response.ConfigFeature, tc.wantFeatureCnt)

			featureMap := make(map[string]*pb.ConfigFeature)
			for _, f := range response.ConfigFeature {
				featureMap[f.Name] = f
			}

			aaaFeature := featureMap["aaa"]
			assert.NotNil(t, aaaFeature)
			assert.Contains(t, aaaFeature.Configuration, tc.wantAAAContains)

			loggerFeature, hasLogger := featureMap["logger"]
			assert.Equal(t, tc.wantHasLogger, hasLogger)
			if tc.wantHasLogger {
				assert.Contains(t, loggerFeature.Configuration, tc.wantLogContains)
			}
		})
	}
}

func TestSetEntityResetsParamsOnSwitchAndNil(t *testing.T) {
	// Save and restore global Registry
	origRegistry := Registry
	defer func() { Registry = origRegistry }()

	Registry = newGeneratorsRegistry()
	err := Registry.RegisterFeatures()
	assert.NoError(t, err)

	request := &pb.ConfigGenRequest{
		Device: &pb.Device{
			Name:   "rt01.foo01",
			Vendor: pb.Vendor_VD_CISCO,
		},
	}
	server := NewServer()

	// 1) Set entity with custom logger params.
	Registry.SetEntity(&entity.ResolvedEntity{
		Name:     "A",
		Features: []string{"aaa", "logger"},
		Params: map[string]map[string]interface{}{
			"logger": {
				"management_interface": "loopback77",
				"logger_server_ips":    []interface{}{"9.9.9.9"},
			},
		},
	})

	respA, err := server.GetConfigGen(context.Background(), request)
	assert.NoError(t, err)
	featuresA := make(map[string]*pb.ConfigFeature)
	for _, f := range respA.ConfigFeature {
		featuresA[f.Name] = f
	}
	assert.Contains(t, featuresA["logger"].Configuration, "loopback77")
	assert.Contains(t, featuresA["logger"].Configuration, "9.9.9.9")

	// 2) Switch to entity with logger enabled but without logger params.
	// Logger must fall back to defaults instead of keeping previous values.
	Registry.SetEntity(&entity.ResolvedEntity{
		Name:     "B",
		Features: []string{"aaa", "logger"},
		Params: map[string]map[string]interface{}{
			"aaa": {"secret": "only-aaa-custom"},
		},
	})

	respB, err := server.GetConfigGen(context.Background(), request)
	assert.NoError(t, err)
	featuresB := make(map[string]*pb.ConfigFeature)
	for _, f := range respB.ConfigFeature {
		featuresB[f.Name] = f
	}
	assert.Contains(t, featuresB["logger"].Configuration, "loopback0")
	assert.Contains(t, featuresB["logger"].Configuration, "192.168.1.1")
	assert.NotContains(t, featuresB["logger"].Configuration, "loopback77")
	assert.NotContains(t, featuresB["logger"].Configuration, "9.9.9.9")

	// 3) Set nil entity. Defaults should still be used.
	Registry.SetEntity(nil)
	respNil, err := server.GetConfigGen(context.Background(), request)
	assert.NoError(t, err)
	featuresNil := make(map[string]*pb.ConfigFeature)
	for _, f := range respNil.ConfigFeature {
		featuresNil[f.Name] = f
	}
	assert.Contains(t, featuresNil["logger"].Configuration, "loopback0")
	assert.Contains(t, featuresNil["logger"].Configuration, "192.168.1.1")
}

func sampleConfigPath(t *testing.T, fileName string) string {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to resolve caller path")
	}
	return filepath.Join(filepath.Dir(thisFile), "..", "configs", "samples", fileName)
}
