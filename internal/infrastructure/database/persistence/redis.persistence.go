package persistence

import (
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
	"gitlab.avakatan.ir/boilerplates/go-boiler/config"
)

func NewRedisClient(connectionConfig config.DatabaseConfig) (*redis.Client, error) {
	fmt.Println(connectionConfig)
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

	return redisClient, nil
}
