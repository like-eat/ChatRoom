package redis

import (
	"apylee_chat_server/pkg/log"
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

var redisClient *redis.Client
var ctx = context.Background()

func createClient() *redis.Client {
	if redisClient == nil {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "dancernumber1",
			DB:       0,
		})
	}
	return redisClient
}

func SetKeyEx(key string, value string, timeout time.Duration) error {
	err := createClient().Set(ctx, key, value, timeout).Err()
	if err != nil {
		log.GetZapLogger().Error(err.Error())
		return err
	}
	return nil
}

func GetKey(key string) (string, error) {
	value, err := createClient().Get(ctx, key).Result()
	if err == redis.Nil {
		log.GetZapLogger().Error("key_with_expire has expired")
		return "", err
	}
	return value, nil
}
