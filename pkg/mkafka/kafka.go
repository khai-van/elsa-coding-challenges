package mkafka

import (
	"time"

	"github.com/segmentio/kafka-go"
)

var sConf Config

type Config struct {
	Brokers []string
}

func InitKafkaConf(conf Config) {
	sConf = conf
}

func GetKafkaWriter(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(sConf.Brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func GetKafkaReader(topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  sConf.Brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 1,    // 1B
		MaxBytes: 10e6, // 10MB
		MaxWait:  100 * time.Millisecond,
	})
}
