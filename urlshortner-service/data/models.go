package data

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type Models struct {
	ShortnerEntry Request
}

type Request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"customShort"`
	ExpireTime  time.Duration `json:"expireTime"`
}

var Context = context.Background()

func CreateClients(No int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis://redis:6379",
		Password: "pass",
		DB:       No,
	})
	return rdb
}
