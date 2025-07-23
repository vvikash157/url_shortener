package config

import (
	"context"
	"log"

	"github.com/go-redis/redis"
)

var (
	redisClient *redis.Client
	Ctx         = context.Background()
)

func InitRedis() *redis.Client {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		DB:       0,
		Password: "",
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Fatalf("error while connecting redis %v", err)
	}
	return redisClient
}
