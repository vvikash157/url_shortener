package repository

import (
	"time"

	"github.com/go-redis/redis"
)

type CacheRepository interface {
	Get(key string) (string, error)
	Set(key, value string, expiration time.Duration) error
}

type RedisCacheRepository struct {
	client *redis.Client
}

func NewRedisCacheClient(client *redis.Client) *RedisCacheRepository {
	return &RedisCacheRepository{client: client}
}

func (r *RedisCacheRepository) Get(key string) (string,error){
	return r.client.Get(key).Result()
}

func (r *RedisCacheRepository) Set(key , value string,expiration time.Duration) error{
	return r.client.Set(key,value,expiration).Err()
}