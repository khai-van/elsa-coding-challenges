package quizserver

import (
	"context"
	pb "quiz/api/gen/quiz"
)

func (srv *Server) SubmitAnswer(ctx context.Context, request *pb.AnswerRequest) (*pb.AnswerResponse, error) {
	return srv.quizSvc.SubmitAnswer(ctx, request.QuizID, request.UserID, request.QuestionID, request.Answer)
}

func (srv *Server) GetLeaderboard(ctx context.Context, request *pb.LeaderboardRequest) (*pb.LeaderboardResponse, error) {
	return srv.leaderboardSvc.GetLeaderboard(ctx, request.QuizID)
}
