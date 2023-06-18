package repository

import (
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/persistence"
)

type RedisRepository struct {
	Client *persistence.RedisClient
}

func NewRedisRepository(client *persistence.RedisClient) *RedisRepository {
	return &RedisRepository{Client: client}
}

func (r *RedisRepository) Get(key string) (string, error) {
	return r.Client.Client.Get(key).Result()
}

func (r *RedisRepository) Set(key string, value string) error {
	return r.Client.Client.Set(key, value, 0).Err()
}

func (r *RedisRepository) Keys(pattern string) ([]string, error) {
	return r.Client.Client.Keys(pattern).Result()
}

// Add more methods as per your requirements
