package encoding

import (
	"encoding/json"
	"fmt"

	"github.com/jaegertracing/jaeger/cmd/opentelemetry/app/exporter/elasticsearchexporter/esmodeltranslator"
	"go.opentelemetry.io/collector/consumer/pdata"
)

type FastJsonMarshal struct {
	Translator *esmodeltranslator.Translator
}

func (f FastJsonMarshal) Marshal(traces pdata.Traces) ([]Message, error) {
	spans, err := f.Translator.ConvertSpans(traces)
	fmt.Println(err, spans)
	var errs []error
	dropped := 0
	messages := make([]Message, 0, len(spans))
	for _, span := range spans {
		data, err := json.Marshal(span.DBSpan)
		fmt.Println(string(data))
		if err != nil {
			errs = append(errs, err)
			dropped++
			continue
		}
		messages = append(messages, Message{Value: data})
	}
	return messages, nil
}

func (f FastJsonMarshal) Encoding() string {
	return "fast_json"
}
