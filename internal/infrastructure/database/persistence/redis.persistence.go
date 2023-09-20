package persistence

import (
	"strconv"

	"github.com/go-redis/redis"
	"github.com/mi6ule/goboy/config"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(connectionConfig config.DatabaseConfig) (*RedisClient, error) {
	dbName, err := strconv.Atoi(connectionConfig.Name)
	if err != nil {
		return nil, err
	}
	// Redis client configuration
	redisClient := redis.NewClient(&redis.Options{
		Addr:     connectionConfig.Host,
		Password: connectionConfig.Pwd, // If Redis requires authentication
		DB:       dbName,               // Redis database index
	})

	return &RedisClient{Client: redisClient}, nil
}
