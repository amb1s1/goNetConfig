package entity

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testConfig = `version: "1.0"

params:
  companyA:
    aaa:
      secret: "companyA-secret-123"
    logger:
      management_interface: "loopback0"
      logger_server_ips:
        - "10.0.1.1"
        - "10.0.1.2"
  companyB:
    aaa:
      secret: "companyB-secret-456"
    logger:
      management_interface: "loopback1"
      logger_server_ips:
        - "172.16.0.1"

features:
  companyA:
    - aaa
    - logger
  companyB:
    - aaa

entity:
  companyA:
    features: companyA
    params: companyA
  companyB:
    features: companyB
    params: companyB

selected_entity: companyA
`

func writeTestConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "entity.yml")
	err := os.WriteFile(path, []byte(content), 0644)
	require.NoError(t, err)
	return path
}

func TestLoad(t *testing.T) {
	path := writeTestConfig(t, testConfig)
	cfg, err := Load(path)
	require.NoError(t, err)
	assert.Equal(t, "1.0", cfg.Version)
	assert.Equal(t, "companyA", cfg.SelectedEntity)
	assert.Contains(t, cfg.Entity, "companyA")
	assert.Contains(t, cfg.Entity, "companyB")
}

func TestLoadFileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/entity.yml")
	assert.Error(t, err)
}

func TestLoadInvalidYAML(t *testing.T) {
	path := writeTestConfig(t, "{{invalid yaml")
	_, err := Load(path)
	assert.Error(t, err)
}

func TestGetSelectedEntityCompanyA(t *testing.T) {
	path := writeTestConfig(t, testConfig)
	cfg, err := Load(path)
	require.NoError(t, err)

	resolved, err := cfg.GetSelectedEntity()
	require.NoError(t, err)

	assert.Equal(t, "companyA", resolved.Name)
	assert.Equal(t, []string{"aaa", "logger"}, resolved.Features)

	// Check AAA params
	aaaParams, ok := resolved.Params["aaa"]
	require.True(t, ok)
	assert.Equal(t, "companyA-secret-123", aaaParams["secret"])

	// Check logger params
	loggerParams, ok := resolved.Params["logger"]
	require.True(t, ok)
	assert.Equal(t, "loopback0", loggerParams["management_interface"])
}

func TestGetSelectedEntityCompanyB(t *testing.T) {
	path := writeTestConfig(t, testConfig)
	cfg, err := Load(path)
	require.NoError(t, err)

	cfg.SelectedEntity = "companyB"
	resolved, err := cfg.GetSelectedEntity()
	require.NoError(t, err)

	assert.Equal(t, "companyB", resolved.Name)
	assert.Equal(t, []string{"aaa"}, resolved.Features)
	assert.Equal(t, "companyB-secret-456", resolved.Params["aaa"]["secret"])
}

func TestGetSelectedEntityNoSelection(t *testing.T) {
	path := writeTestConfig(t, testConfig)
	cfg, err := Load(path)
	require.NoError(t, err)

	cfg.SelectedEntity = ""
	_, err = cfg.GetSelectedEntity()
	assert.Error(t, err)
}

func TestGetSelectedEntityInvalidName(t *testing.T) {
	path := writeTestConfig(t, testConfig)
	cfg, err := Load(path)
	require.NoError(t, err)

	cfg.SelectedEntity = "nonexistent"
	_, err = cfg.GetSelectedEntity()
	assert.Error(t, err)
}

func TestSetEntity(t *testing.T) {
	path := writeTestConfig(t, testConfig)

	err := SetEntity(path, "companyB")
	require.NoError(t, err)

	// Reload and verify
	cfg, err := Load(path)
	require.NoError(t, err)
	assert.Equal(t, "companyB", cfg.SelectedEntity)
}

func TestSetEntityInvalidName(t *testing.T) {
	path := writeTestConfig(t, testConfig)
	err := SetEntity(path, "nonexistent")
	assert.Error(t, err)
}
