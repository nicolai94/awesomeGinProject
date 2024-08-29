package utils

import (
	"awesomeProject/config"
	"context"
	"time"
)

func AddToRedis(tokenString string) error {
	ctx := context.Background()
	err := config.RedisClient.Set(ctx, tokenString, true, 24*time.Hour).Err()
	return err
}

// RemoveFromRedis удаляет токен из Redis
func RemoveFromRedis(tokenString string) error {
	ctx := context.Background()
	err := config.RedisClient.Del(ctx, tokenString).Err()
	return err
}
