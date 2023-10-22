package initializers

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gym_friend_auth_server/utils"
	"time"
)

var InMemoryDB InMemoryDatabase

type InMemoryDatabase interface {
	Set(key string, value string) error
	SetExp(key string, value string, expires time.Duration) error
	Get(key string) (string, error)
	Del(key string) error
}

type Redis struct {
	client *redis.Client
}

// redis server 연결을 위한 코드입니다
// 데이터 캐싱을 위해 사용하는 inmemory DB입니다
func InMemoryConnection() {
	r := &Redis{}
	r.client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	InMemoryDB = r
}

func (r *Redis) Set(key string, value string) error {
	return r.client.Set(context.Background(), key, value, 0).Err()
}

func (r *Redis) SetExp(key string, value string, expires time.Duration) error {
	return r.client.Set(context.Background(), key, value, expires).Err()
}

func (r *Redis) Get(key string) (string, error) {
	str, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return "", utils.ErrValueNotFound
	}
	return str, err
}

func (r *Redis) Del(key string) error {
	return r.client.Del(context.Background(), key).Err()
}
