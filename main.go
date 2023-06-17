package main

import (
	"log"

	"github.com/amb1s1/gonetconfig/server"
	"github.com/amb1s1/gonetconfig/service"
)

func main() {
	if err := service.Registry.RegisterFeatures(); err != nil {
		log.Fatalf("RegisterFeatures(), failed to register features: %v", err)
		return
	}
	if err := server.Start(); err != nil {
		log.Fatalf("start(), failed to start the server: %v", err)
		return
	}
}
