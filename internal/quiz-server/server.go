package quizserver

import (
	"context"
	pb "quiz/api/gen/quiz"
	"quiz/internal/modules/leaderboard"
	quizservice "quiz/internal/modules/quiz-service"
	"quiz/pkg/config"
	"quiz/pkg/mkafka"
	"quiz/pkg/mredis"
	"time"
)

type Server struct {
	config Config

	pb.UnimplementedQuizServiceServer

	// service
	quizSvc        *quizservice.Service
	leaderboardSvc *leaderboard.Service
}

func New() (*Server, error) {
	conf, err := config.LoadConfig[Config]("quiz")
	if err != nil {
		return nil, err
	}

	server := Server{
		config: conf,
	}

	// init connection
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5) // 5 second timeout
	defer cancel()

	// connect redis
	if err := mredis.ConnectRedis(ctx, conf.Redis); err != nil {
		return nil, err
	}
	// connect kafka
	mkafka.InitKafkaConf(conf.Kafka)

	// new service handler
	server.leaderboardSvc = leaderboard.New()
	server.quizSvc = quizservice.New(server.leaderboardSvc)

	return &server, nil
}
