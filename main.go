package main

import (
	"log"

	"github.com/amb1s1/gonetconfig/entity"
	"github.com/amb1s1/gonetconfig/server"
	"github.com/amb1s1/gonetconfig/service"
)

func main() {
	if err := service.Registry.RegisterFeatures(); err != nil {
		log.Fatalf("RegisterFeatures(), failed to register features: %v", err)
		return
	}

	// Load entity config: try user config first, fall back to template.
	entityCfg, err := loadEntityConfig()
	if err != nil {
		log.Printf("entity config not loaded, using defaults: %v", err)
	} else {
		resolved, err := entityCfg.GetSelectedEntity()
		if err != nil {
			log.Printf("failed to resolve entity, using defaults: %v", err)
		} else {
			service.Registry.SetEntity(resolved)
			log.Printf("entity loaded: %s", resolved.Name)
		}
	}

	if err := server.Start(); err != nil {
		log.Fatalf("start(), failed to start the server: %v", err)
		return
	}
}

// loadEntityConfig tries the user config path first, then falls back to the
// template config shipped with the repo.
func loadEntityConfig() (*entity.Config, error) {
	defaultPath, err := entity.DefaultPath()
	if err == nil {
		cfg, err := entity.Load(defaultPath)
		if err == nil {
			return cfg, nil
		}
	}
	return entity.Load(entity.TemplateConfigPath)
}
