package kafka

import (
	"context"
	"encoding/json"
	"flowanalysis/pkg/log"

	"github.com/segmentio/kafka-go"
)

var (
	kafkaStream *KafkaWriter
)

const (
	broker            = "kafka:9092"
	FLOW_UNIQUE_TOPIC = "flow.unique.entries"
)

type KafkaWriter struct {
	writer *kafka.Writer
}

func InitKafka() {
	kafkaStream = NewKafkaWriter(FLOW_UNIQUE_TOPIC, broker)
}

func NewKafkaWriter(topic, broker string) *KafkaWriter {
	return &KafkaWriter{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:  []string{broker},
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		}),
	}
}

func (k *KafkaWriter) SendMessage(ctx context.Context, message interface{}) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = k.writer.WriteMessages(ctx, kafka.Message{
		Value: jsonData,
	})
	if err != nil {
		log.Print(log.Error, "Failed to send message to Kafka: %v", err)
		return err
	}

	log.Print(log.Info, "Message sent to Kafka: %s", string(jsonData))
	return nil
}

func (k *KafkaWriter) Close() error {
	return k.writer.Close()
}

func SendMessage(ctx context.Context, topic, message interface{}) {
	switch topic {
	case FLOW_UNIQUE_TOPIC:
		kafkaStream.SendMessage(ctx, message)
	default:
		log.Print(log.Error, "Not a valid topic")
	}
}
