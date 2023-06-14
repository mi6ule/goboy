package repository

import "github.com/go-redis/redis"

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{client: client}
}

func (r *RedisRepository) Get(key string) (string, error) {
	return r.client.Get(key).Result()
}

func (r *RedisRepository) Set(key string, value string) error {
	return r.client.Set(key, value, 0).Err()
}

func (r *RedisRepository) Keys(pattern string) ([]string, error) {
	return r.client.Keys(pattern).Result()
}

// Add more methods as per your requirements
