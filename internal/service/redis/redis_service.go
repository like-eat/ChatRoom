package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"kama_chat_server/config"
	"kama_chat_server/pkg/zlog"
	"strconv"
	"time"
)

var redisClient *redis.Client
var ctx = context.Background()

func init() {
	conf := config.GetConfig()
	host := conf.RedisConfig.Host
	port := conf.RedisConfig.Port
	password := conf.RedisConfig.Password
	db := conf.Db
	addr := host + ":" + strconv.Itoa(port)

	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

func SetKeyEx(key string, value string, timeout time.Duration) error {
	err := redisClient.Set(ctx, key, value, timeout).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetKey(key string) (string, error) {
	value, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			zlog.Info("该key不存在")
			return "", nil
		}
		return "", err
	}
	return value, nil
}

func DelKey(key string) error {
	err := redisClient.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
