package mkafka

import (
	"github.com/IBM/sarama"
)

// ConsumeMessages return a channel messages from a specific topic
func ConsumeMessages(topic string) (sarama.PartitionConsumer, error) {
	consumer, err := GetKafkaConsumer()
	if err != nil {
		return nil, err
	}

	// Partition consumer
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		return nil, err
	}

	return partitionConsumer, nil
}
