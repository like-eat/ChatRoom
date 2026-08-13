package redis

import (
	"context"
	"kama_chat_server/internal/config"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

// 上下文
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

// 存一个带过期时间的键值对，主要负责查到数据库后把数据塞入redis
func SetKeyEx(key string, value string, timeout time.Duration) error {
	// key是存哪个键值对
	// value是存哪个值
	// timeout是过期时间
	err := redisClient.Set(ctx, key, value, timeout).Err()
	if err != nil {
		return err
	}
	return nil
}

// 如果key不存在则返回错误
func GetKeyNilIsErr(key string) (string, error) {
	value, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

// ScanAndDelete 使用 SCAN 命令批量删除匹配的 key（生产推荐，不会阻塞 Redis）
// 一次只取一部分，删完再取下一批，不会阻塞redis
func ScanAndDelete(pattern string) error {
	var cursor uint64 = 0
	for {
		keys, nextCursor, err := redisClient.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if _, err := redisClient.Del(ctx, keys...).Result(); err != nil {
				return err
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}

// 全量扫描，直到全部删除完成
func DelKeysWithPattern(pattern string) error {
	var keys []string
	var err error

	for {
		// 使用 Keys 命令迭代匹配的键
		keys, err = redisClient.Keys(ctx, pattern).Result()
		if err != nil {
			return err
		}

		// 如果没有更多的键，则跳出循环
		if len(keys) == 0 {
			log.Println("没有找到对应key")
			break
		}

		// 删除找到的键
		if len(keys) > 0 {
			_, err = redisClient.Del(ctx, keys...).Result()
			if err != nil {
				return err
			}
			log.Println("成功删除相关对应key", keys)
		}
	}

	return nil
}

// 删除所有key
func DeleteAllRedisKeys() error {
	var cursor uint64 = 0
	for {
		keys, nextCursor, err := redisClient.Scan(ctx, cursor, "*", 0).Result()
		if err != nil {
			return err
		}
		cursor = nextCursor

		if len(keys) > 0 {
			_, err := redisClient.Del(ctx, keys...).Result()
			if err != nil {
				return err
			}
		}

		if cursor == 0 {
			break
		}
	}
	return nil
}
