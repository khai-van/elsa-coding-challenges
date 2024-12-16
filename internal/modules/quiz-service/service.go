package quizservice

import (
	"context"
	"math/rand/v2"
	pb "quiz/api/gen/quiz"
)

type LeaderboardService interface {
	AddScore(ctx context.Context, quizID, userID string, score int32) error
}

type Service struct {
	lbSvc LeaderboardService
}

func New(lbSvc LeaderboardService) *Service {
	return &Service{
		lbSvc: lbSvc,
	}
}

func (svc *Service) SubmitAnswer(ctx context.Context, quizID, userID, questionID, answer string) (*pb.AnswerResponse, error) {
	// TODO: validate data

	// add random score to  lb
	score := int32(rand.Uint32N(20))
	if err := svc.lbSvc.AddScore(ctx, quizID, userID, score); err != nil {
		return nil, err
	}

	return &pb.AnswerResponse{
		IsCorrect: true,
		Score:     score,
	}, nil
}
