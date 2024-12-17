package leaderboard

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	pb "quiz/api/gen/quiz"
	"quiz/internal/constant"
	"quiz/models"
	"quiz/pkg/mkafka"
	"quiz/pkg/mredis"
	"time"

	"github.com/segmentio/kafka-go"
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
	result, err := cli.ZRevRangeWithScores(ctx, svc.leaderboarKey(quizID), 0, -1).Result() // get all
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
	err := cli.ZIncrBy(ctx, key, float64(score), userID).Err()
	if err != nil {
		return err
	}

	// publish member change to kafka
	go svc.publishChangeLeaderboard(key, quizID, userID)

	return nil
}

func (*Service) publishChangeLeaderboard(key string, quizID, userID string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	cli := mredis.GetClient()
	// get new rank
	rank, err := cli.ZRevRankWithScore(ctx, key, userID).Result()
	if err != nil {
		log.Println(err)
		return
	}

	msg := models.LeaderboardMemberChange{
		QuizID:   quizID,
		UserID:   userID,
		NewRank:  int(rank.Rank),
		NewScore: int(rank.Score),
	}
	msgStr, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return
	}

	// publish msg
	if err := mkafka.GetKafkaWriter(constant.KAFKATOPIC_LEADERBOARDCHANGE).WriteMessages(
		ctx,
		kafka.Message{
			Key:   []byte(fmt.Sprintf("%s:%s", key, userID)),
			Value: msgStr,
		},
	); err != nil {
		log.Println(err)
	}
}
