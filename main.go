package main

import (
	"log"

	"github.com/amb1s1/gonetconfig/service"
)

func main() {
	if err := service.Registry.RegisterFeatures(); err != nil {
		log.Fatalln(err)
	}
	if err := service.Start(); err != nil {
		log.Fatalln(err)
	}
}
