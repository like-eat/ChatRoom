package redis

import (
	"apylee_chat_server/config"
	"apylee_chat_server/pkg/log"
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

var redisClient *redis.Client
var ctx = context.Background()

func createClient() *redis.Client {
	if redisClient == nil {
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
