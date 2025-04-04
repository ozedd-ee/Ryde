package cache

import (
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

// Initialize Redis client
func Init() *redis.Client {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL environment variable not set")
	}
	return redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
}
