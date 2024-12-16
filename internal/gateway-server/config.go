package gatewayserver

import "quiz/pkg/mkafka"

type Config struct {
	Kafka        mkafka.Config
	QuizInternal string 
}
