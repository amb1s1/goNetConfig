package entity

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	// DefaultConfigDir is the directory under the user's home for gonetconfig.
	DefaultConfigDir = ".gonetconfig"
	// DefaultConfigFile is the entity config filename.
	DefaultConfigFile = "entity.yml"
	// TemplateConfigPath is the path to the template config shipped with the repo.
	TemplateConfigPath = "configs/entity.yml"
)

// Config represents the parsed entity YAML configuration.
type Config struct {
	Version        string                                   `yaml:"version"`
	Params         map[string]map[string]map[string]interface{} `yaml:"params"`
	Features       map[string][]string                      `yaml:"features"`
	Entity         map[string]EntityDef                     `yaml:"entity"`
	SelectedEntity string                                   `yaml:"selected_entity"`
}

// EntityDef maps an entity name to its features and params keys.
type EntityDef struct {
	Features string `yaml:"features"`
	Params   string `yaml:"params"`
}

// ResolvedEntity holds the resolved features and params for the selected entity.
type ResolvedEntity struct {
	Name     string
	Features []string
	Params   map[string]map[string]interface{}
}

// Load parses an entity YAML file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read entity config %s: %w", path, err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse entity config %s: %w", path, err)
	}
	return &cfg, nil
}

// DefaultPath returns the default entity config path (~/.gonetconfig/entity.yml).
func DefaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, DefaultConfigDir, DefaultConfigFile), nil
}

// GetSelectedEntity resolves the active entity's features and params.
func (c *Config) GetSelectedEntity() (*ResolvedEntity, error) {
	if c.SelectedEntity == "" {
		return nil, fmt.Errorf("no selected_entity defined in config")
	}
	entDef, ok := c.Entity[c.SelectedEntity]
	if !ok {
		return nil, fmt.Errorf("entity %q not found in entity definitions", c.SelectedEntity)
	}

	features, ok := c.Features[entDef.Features]
	if !ok {
		return nil, fmt.Errorf("features key %q not found for entity %q", entDef.Features, c.SelectedEntity)
	}

	params, ok := c.Params[entDef.Params]
	if !ok {
		return nil, fmt.Errorf("params key %q not found for entity %q", entDef.Params, c.SelectedEntity)
	}

	return &ResolvedEntity{
		Name:     c.SelectedEntity,
		Features: features,
		Params:   params,
	}, nil
}

// SetEntity updates the selected_entity field in the YAML file at the given path.
func SetEntity(path, name string) error {
	cfg, err := Load(path)
	if err != nil {
		return err
	}
	if _, ok := cfg.Entity[name]; !ok {
		return fmt.Errorf("entity %q not found in config", name)
	}
	cfg.SelectedEntity = name

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal entity config: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write entity config %s: %w", path, err)
	}
	return nil
}
