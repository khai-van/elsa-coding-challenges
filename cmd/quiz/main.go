package main

import (
	"log"
	"net"

	pb "quiz/api/gen/quiz"
	quizserver "quiz/internal/quiz-server"

	"google.golang.org/grpc"
)

func main() {
	server := grpc.NewServer()
	quizHandler, err := quizserver.New()
	if err != nil {
		log.Fatalf("Failed to create new server: %v", err)
	}
	pb.RegisterQuizServiceServer(server, quizHandler)

	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
