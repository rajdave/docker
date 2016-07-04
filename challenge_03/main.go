package main

import (
	"log"

	"github.com/docker/go-plugins-helpers/volume"
)

func main() {
	driver := NewRDriver()
	handler := volume.NewHandler(driver)
	if err := handler.ServeUnix("root", "RDriver-example"); err != nil {
		log.Fatalf("Error %v", err)
	}

	for {

	}
}
