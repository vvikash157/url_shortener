package config

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis"
)

var (
	redisClient *redis.Client
	Ctx         = context.Background()
)

func InitRedis() *redis.Client {
	dbStr := os.Getenv("REDIS_DB")
	db := 0 // default DB
	if dbStr != "" {
		var err error
		db, err = strconv.Atoi(dbStr)
		if err != nil {
			log.Fatalf("Invalid REDIS_DB: %v", err)
		}
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Fatalf("error while connecting redis %v", err)
	}
	return redisClient
}
