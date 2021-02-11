package main

import (
	"github.com/open-telemetry/graylog"
	"log"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/service"
)

func main() {
	factories, err := graylog.Components()
	if err != nil {
		log.Fatalf("failed to build components: %v", err)
	}

	info := component.ApplicationStartInfo{
		ExeName:  "graylog-custom",
		LongName: "gray custom exporter",
		Version:  "1.0.0",
	}

	app, err := service.New(service.Parameters{
		Factories:            factories,
		ApplicationStartInfo: info,
	})

	if err != nil {
		log.Fatalf("failed to construct application: %v", err)
	}

	err = app.Run()
	if err != nil {
		log.Fatalf("application run finished with error: %v", err)
	}
}
