package graylog

import (
	"github.com/open-telemetry/graylog/kafka"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenterror"
	"go.opentelemetry.io/collector/service/defaultcomponents"
)

func Components() (component.Factories, error) {
	var errs []error
	factories, err := defaultcomponents.Components()
	if err != nil {
		return component.Factories{}, err
	}

	var exporters []component.ExporterFactory

	for _, ep := range factories.Exporters {
		exporters = append(exporters, ep)
	}

	exporters = append(exporters, kafka.NewFactory())
	factories.Exporters, err = component.MakeExporterFactoryMap(exporters...)
	if err != nil {
		errs = append(errs, err)
	}
	return factories, componenterror.CombineErrors(errs)
}
