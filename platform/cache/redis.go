package cache

import (
	// "fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

// RedisConnection func for connect to Redis server.
func RedisConnection() (*redis.Client, error) {
	// Define Redis database number.
	dbNumber, _ := strconv.Atoi(os.Getenv("REDIS_DB_NUMBER"))

	// Set default Redis host if REDIS_HOST is empty.
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost" // or "127.0.0.1"
	}

	// // URL for Redis connection.
	// redisConnURL := fmt.Sprintf(
	// 	"%s:%s",
	// 	redisHost,
	// 	os.Getenv("REDIS_PORT"),
	// )

	// Set Redis options.
	options := &redis.Options{
		Addr: "cache:6379",
		// Addr:     redisConnURL,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbNumber,
	}

	return redis.NewClient(options), nil
}
