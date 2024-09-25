// redis.go
package database

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client
var ctx = context.Background()
var RedisClient *redis.Client

// ConnectRedis menghubungkan ke Redis
func ConnectRedis() *redis.Client {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Could not connect to Redis: %v", err))
	}

	fmt.Println("Connected to Redis!")
	return RedisClient
}

func GetRedisClient() *redis.Client {
	return rdb
}

// SetKey sets a key-value pair in Redis
func SetKey(ctx context.Context, key string, value interface{}) error {
	err := RedisClient.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetKey gets the value for a key from Redis
func GetKey(ctx context.Context, key string) (string, error) {
	val, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
