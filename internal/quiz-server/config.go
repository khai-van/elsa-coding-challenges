package quizserver

import (
	"quiz/pkg/mkafka"
	"quiz/pkg/mredis"
)

type Config struct {
	Redis mredis.Config
	Kafka mkafka.Config
}
