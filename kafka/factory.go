package kafka

import (
	"context"
	"time"

	"github.com/open-telemetry/graylog/config"

	"github.com/open-telemetry/graylog/encoding"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	typeStr                     = "graylogkafka"
	defaultTracesTopic          = "otlp_spans"
	defaultEncoding             = "fast_json"
	defaultBroker               = "localhost:9092"
	defaultMetadataFull         = true
	defaultMetadataRetryMax     = 3
	defaultMetadataRetryBackoff = time.Millisecond * 250
)

type kafkaExporterFactory struct {
	traceMarshallers map[string]encoding.TraceMarshaller
}

type FactoryOption func(factory *kafkaExporterFactory)

func NewFactory(options ...FactoryOption) component.ExporterFactory {
	f := &kafkaExporterFactory{}

	for _, o := range options {
		o(f)
	}
	return exporterhelper.NewFactory(
		typeStr,
		createDefaultConfig,
		exporterhelper.WithTraces(f.createTraceExporter),
	)
}

func createDefaultConfig() configmodels.Exporter {
	return &config.Config{
		ExporterSettings: configmodels.ExporterSettings{
			TypeVal: typeStr,
			NameVal: typeStr,
		},
		TimeoutSettings: exporterhelper.DefaultTimeoutSettings(),
		RetrySettings:   exporterhelper.DefaultRetrySettings(),
		QueueSettings:   exporterhelper.DefaultQueueSettings(),
		Brokers:         []string{defaultBroker},

		Topic:    "",
		Encoding: defaultEncoding,
		Metadata: config.Metadata{
			Full: defaultMetadataFull,
			Retry: config.MetadataRetry{
				Max:     defaultMetadataRetryMax,
				Backoff: defaultMetadataRetryBackoff,
			},
		},
	}
}

func (f *kafkaExporterFactory) createTraceExporter(_ context.Context, params component.ExporterCreateParams, cfg configmodels.Exporter) (component.TracesExporter, error) {
	oCfg := cfg.(*config.Config)
	if oCfg.Topic == "" {
		oCfg.Topic = defaultTracesTopic
	}
	f.traceMarshallers = encoding.TracesMarshallers(oCfg)
	exp, err := newTracerExporterProducer(*oCfg, params, f.traceMarshallers)
	if err != nil {
		return nil, err
	}
	return exporterhelper.NewTraceExporter(cfg,
		params.Logger, exp.traceDataPusher,
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 0}),
		exporterhelper.WithRetry(oCfg.RetrySettings),
		exporterhelper.WithQueue(oCfg.QueueSettings),
		exporterhelper.WithShutdown(exp.Close))
	return nil, nil
}
