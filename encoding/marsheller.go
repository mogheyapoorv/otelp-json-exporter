package encoding

import (
	"github.com/jaegertracing/jaeger/cmd/opentelemetry/app/exporter/elasticsearchexporter/esmodeltranslator"
	"github.com/open-telemetry/graylog/config"
	"go.opentelemetry.io/collector/consumer/pdata"
)

type TraceMarshaller interface {
	Marshal(traces pdata.Traces) ([]Message, error)
	Encoding() string
}

type Message struct {
	Value []byte
}

// tracesMarshallers returns map of supported encodings with TracesMarshaller.
func TracesMarshallers(params *config.Config) map[string]TraceMarshaller {
	tagsKeysAsFields, _ := params.TagKeysAsFields()

	fastJson := &FastJsonMarshal{
		Translator: esmodeltranslator.NewTranslator(params.AllAsFields, tagsKeysAsFields, params.DotReplacement),
	}
	return map[string]TraceMarshaller{
		fastJson.Encoding(): fastJson,
	}
}
