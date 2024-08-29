package config

import (
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
)

var RedisClient *redis.Client

func ConnectToDB() *gorm.DB {
	var err error
	dsn := os.Getenv("DB_DSN")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database. Error: ", err)
	}

	return db
}

func ConnectRedis() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	redisDBstr := os.Getenv("REDIS_DB")
	redisDB, err := strconv.Atoi(redisDBstr)
	if err != nil {
		log.Fatalf("Error converting REDIS_DB to integer: %v", err)
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   redisDB,
	})
}
