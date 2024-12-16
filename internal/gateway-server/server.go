package gatewayserver

import (
	"quiz/pkg/config"
	"quiz/pkg/mkafka"

	pb "quiz/api/gen/quiz"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	config Config
	hub    Hub

	quizService pb.QuizServiceClient
}

func New() (*Server, error) {
	conf, err := config.LoadConfig[Config]("gateway")
	if err != nil {
		return nil, err
	}

	server := Server{
		config: conf,
		hub: Hub{
			Users: map[string]*websocket.Conn{},
			Rooms: map[string]map[string]struct{}{},
		},
	}

	// init connection
	// connect kafka
	if err := mkafka.InitKafka(conf.Kafka); err != nil {
		return nil, err
	}

	// set connection to quiz service
	conn, err := grpc.Dial(conf.QuizInternal, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	server.quizService = pb.NewQuizServiceClient(conn)

	// run hub
	go server.hub.run()

	return &server, nil
}
