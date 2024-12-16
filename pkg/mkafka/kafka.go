package mkafka

import (
	"sync"

	"github.com/IBM/sarama"
)

var (
	kafkaProducerInstance sarama.SyncProducer
	kafkaConsumerInstance sarama.Consumer
	onceProducer          sync.Once
	onceConsumer          sync.Once
	sConf                 Config
)

type Config struct {
	Brokers []string
}

func InitKafka(conf Config) error {
	sConf = conf

	// alway init producer
	if _, err := GetKafkaProducer(); err != nil {
		return err
	}

	return nil
}

func GetKafkaProducer() (sarama.SyncProducer, error) {
	var err error
	onceProducer.Do(func() {
		// Configure Kafka producer
		producerConfig := sarama.NewConfig()
		producerConfig.Producer.Return.Successes = true
		producerConfig.Producer.Retry.Max = 5

		producer, initErr := sarama.NewSyncProducer(sConf.Brokers, producerConfig)
		if initErr != nil {
			err = initErr
			return
		}
		kafkaProducerInstance = producer
	})
	return kafkaProducerInstance, err
}

func GetKafkaConsumer() (sarama.Consumer, error) {
	var err error
	onceConsumer.Do(func() {
		// Configure Kafka consumer
		consumer, initErr := sarama.NewConsumer(sConf.Brokers, nil)
		if initErr != nil {
			err = initErr
			return
		}
		kafkaConsumerInstance = consumer
	})
	return kafkaConsumerInstance, err
}
