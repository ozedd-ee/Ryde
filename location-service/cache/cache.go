package cache

import "github.com/redis/go-redis/v9"

// Initialize Redis client
func Init() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}
