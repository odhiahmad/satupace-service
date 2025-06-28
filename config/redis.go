package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var (
	redisCtx = context.Background()
)

// SetupRedisClient mengembalikan instance *redis.Client
func SetupRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	if _, err := client.Ping(redisCtx).Result(); err != nil {
		log.Fatalf("❌ Gagal koneksi Redis: %v", err)
	}

	log.Println("✅ Redis client berhasil terhubung")
	return client
}
