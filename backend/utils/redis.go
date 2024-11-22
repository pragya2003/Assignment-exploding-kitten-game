package utils

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()
var RedisClient *redis.Client

// InitRedis initializes a connection to the Redis server
func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),     // Redis server address
		Password: os.Getenv("REDIS_PASSWORD"), // Redis password (if any)
		DB:       0,                           // Default DB
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	log.Println("Connected to Redis")
}
