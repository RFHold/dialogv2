package message

import (
	"dialogv2/internal/config"
	"github.com/go-redis/redis/v8"
)

const MutateChannel = "MUTATE_MESSAGE_CHANNEL"

func NewRedisClient(cfg *config.RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Host + ":" + cfg.Port,
		DB:   0,
	})

	return rdb
}
