package mredis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Address  string
	UserName string
	Password string
}

var client redis.UniversalClient

func GetClient() redis.UniversalClient {
	return client
}

func ConnectRedis(ctx context.Context, conf Config) error {
	client = redis.NewClient(&redis.Options{
		Addr:     conf.Address,
		Username: conf.UserName,
		Password: conf.Password,
	})

	if client == nil {
		return fmt.Errorf("can't new client")
	}

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("can't ping to Redis: %w", err)
	}

	return nil
}
