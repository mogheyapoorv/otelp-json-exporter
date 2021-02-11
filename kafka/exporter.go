package kafka

import (
	"context"
	"fmt"

	"github.com/open-telemetry/graylog/config"

	"go.opentelemetry.io/collector/component"

	"github.com/Shopify/sarama"
	"github.com/open-telemetry/graylog/encoding"
	"go.opentelemetry.io/collector/consumer/consumererror"
	"go.opentelemetry.io/collector/consumer/pdata"
	"go.uber.org/zap"
)

type kafkaTracesProducer struct {
	producer   sarama.SyncProducer
	topic      string
	marshaller encoding.TraceMarshaller
	logger     *zap.Logger
}

func (e *kafkaTracesProducer) traceDataPusher(_ context.Context, td pdata.Traces) (int, error) {
	messages, err := e.marshaller.Marshal(td)
	fmt.Println(err)
	if err != nil {
		return td.SpanCount(), consumererror.Permanent(err)
	}

	err = e.producer.SendMessages(producerMessage(messages, e.topic))
	if err != nil {
		return td.SpanCount(), err
	}
	return 0, nil
}

func (e *kafkaTracesProducer) Close(context.Context) error {
	return e.producer.Close()
}

func producerMessage(messages []encoding.Message, topic string) []*sarama.ProducerMessage {
	producerMessages := make([]*sarama.ProducerMessage, len(messages))
	for i := range messages {
		producerMessages[i] = &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(messages[i].Value),
		}
	}
	return producerMessages
}

func newTracerExporterProducer(config config.Config, params component.ExporterCreateParams, marshallers map[string]encoding.TraceMarshaller) (*kafkaTracesProducer, error) {
	marshaller := marshallers[config.Encoding]
	if marshaller == nil {
		return nil, fmt.Errorf("unrecoznige encoding")
	}
	producer, err := newSaramaProducer(config)
	if err != nil {
		return nil, err
	}

	return &kafkaTracesProducer{
		producer:   producer,
		topic:      config.Topic,
		marshaller: marshaller,
		logger:     params.Logger,
	}, nil
}

func newSaramaProducer(config config.Config) (sarama.SyncProducer, error) {
	c := sarama.NewConfig()
	c.Producer.Return.Successes = true
	c.Producer.Return.Errors = true
	c.Producer.RequiredAcks = sarama.WaitForLocal
	c.Producer.Timeout = config.Timeout
	c.Metadata.Full = config.Metadata.Full
	c.Metadata.Retry.Max = config.Metadata.Retry.Max
	c.Metadata.Retry.Backoff = config.Metadata.Retry.Backoff

	if config.ProtocolVersion != "" {
		version, err := sarama.ParseKafkaVersion(config.ProtocolVersion)
		if err != nil {
			return nil, err
		}
		c.Version = version
	}

	producer, err := sarama.NewSyncProducer(config.Brokers, c)
	if err != nil {
		return nil, err
	}
	return producer, nil
}
