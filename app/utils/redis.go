package utils

import (
	"context"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

func ConnectRedis() *redis.Client {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisDBstr := os.Getenv("REDIS_DB")
	redisDB, err := strconv.Atoi(redisDBstr)
	if err != nil {
		log.Fatalf("Error converting REDIS_DB to integer: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: redisAddr, // Адрес Redis-сервера
		DB:   redisDB,   // Номер базы данных Redis
	})

	// Проверка соединения
	ctx := context.Background()
	_, err = client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	} else {
		log.Println("Successfully connected to Redis")
	}

	return client
}

var RedisClient *redis.Client

func ConnectToRedis() {
	RedisClient = ConnectRedis() // Инициализация клиента Redis
}

func AddToRedis(key string, value string) error {
	ctx := context.Background()
	err := RedisClient.Set(ctx, key, value, 24*time.Hour).Err()
	return err
}

func GetFromRedis(key string) (string, error) {
	ctx := context.Background()
	val, err := RedisClient.Get(ctx, key).Result()
	return val, err
}

func RemoveFromRedis(key string) error {
	ctx := context.Background()
	err := RedisClient.Del(ctx, key).Err()
	return err
}
