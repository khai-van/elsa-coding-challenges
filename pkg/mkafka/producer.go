package mkafka

import (
	"github.com/IBM/sarama"
)

// SendMessage sends a message to the given Kafka topic
func SendMessage(topic, message string) error {
	producer, err := GetKafkaProducer()
	if err != nil {
		return err
	}

	// Prepare Kafka message
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	// Send the message
	_, _, err = producer.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}
