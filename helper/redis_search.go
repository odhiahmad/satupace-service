package helper

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// SetJSONToRedis serializes data to JSON and stores it in Redis with a TTL
func SetJSONToRedis(ctx context.Context, rdb *redis.Client, key string, data interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, key, jsonData, ttl).Err()
}

// GetJSONFromRedis retrieves data from Redis and deserializes it from JSON
func GetJSONFromRedis(ctx context.Context, rdb *redis.Client, key string, dest interface{}) error {
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}
