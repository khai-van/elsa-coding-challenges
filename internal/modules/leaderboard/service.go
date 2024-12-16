package leaderboard

import (
	"context"
	"encoding/json"
	"fmt"
	pb "quiz/api/gen/quiz"
	"quiz/internal/constant"
	"quiz/models"
	"quiz/pkg/mkafka"
	"quiz/pkg/mredis"

	"github.com/redis/go-redis/v9"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (svc *Service) leaderboarKey(quizID string) string {
	const LEADERBOARD_PREX_KEY = "Quiz:Leaderboard"

	return fmt.Sprintf("%s:%s", LEADERBOARD_PREX_KEY, quizID)
}

func (svc *Service) GetLeaderboard(ctx context.Context, quizID string) (*pb.LeaderboardResponse, error) {
	// TODO: validate quizID in quiz repo

	// assume that quizID is correct, improve in future using pagination for performance
	// get leaderboard from sorted set redis
	cli := mredis.GetClient()
	result, err := cli.ZRangeWithScores(ctx, svc.leaderboarKey(quizID), 0, -1).Result() // get all
	if err != nil {
		return nil, err
	}

	var resp pb.LeaderboardResponse
	for _, item := range result {
		resp.Leaderboard = append(resp.Leaderboard, &pb.UserScore{
			UserID: item.Member.(string),
			Score:  int32(item.Score),
		})
	}

	return &resp, nil
}

func (svc *Service) AddScore(ctx context.Context, quizID, userID string, score int32) error {
	// assume that already validate output
	// add score to sorted set
	cli := mredis.GetClient()
	key := svc.leaderboarKey(quizID)
	newScore, err := cli.ZAdd(ctx, key, redis.Z{
		Score:  float64(score),
		Member: userID,
	}).Result()
	if err != nil {
		return err
	}

	// publish member change to kafka
	// get new rank
	rank, err := cli.ZRank(ctx, key, userID).Result()
	if err != nil {
		return err
	}

	msg := models.LeaderboardMemberChange{
		UserID:   userID,
		NewRank:  int(rank),
		NewScore: int(newScore),
	}
	msgStr, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// publish msg
	if err := mkafka.SendMessage(constant.KAFKATOPIC_LEADERBOARDCHANGE, string(msgStr)); err != nil {
		return err
	}

	return nil
}
