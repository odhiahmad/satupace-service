// file: helper/redis_helper.go

package helper

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisHelper struct {
	Client *redis.Client
	Ctx    context.Context
}

func NewRedisHelper(client *redis.Client) *RedisHelper {
	return &RedisHelper{
		Client: client,
		Ctx:    context.Background(),
	}
}

func (r *RedisHelper) AllowRequest(identifier string, maxTry int, window time.Duration) error {
	key := fmt.Sprintf("rate_limit:otp:%s", identifier)

	count, err := r.Client.Incr(r.Ctx, key).Result()
	if err != nil {
		return fmt.Errorf("gagal mengecek rate limit OTP: %v", err)
	}

	if count == 1 {
		r.Client.Expire(r.Ctx, key, window)
	}

	if int(count) > maxTry {
		return fmt.Errorf("Terlalu sering meminta OTP. Silakan coba lagi nanti.")
	}

	return nil
}

func (r *RedisHelper) SaveOTP(prefix, identifier, otp string, expiration time.Duration) error {
	key := fmt.Sprintf("otp:%s:%s", prefix, identifier)
	err := r.Client.Set(r.Ctx, key, otp, expiration).Err()
	if err != nil {
		log.Printf("❌ Gagal menyimpan OTP ke Redis: %v", err)
	} else {
		log.Printf("✅ OTP berhasil disimpan [%s] = %s", key, otp)
	}
	return err
}

func (r *RedisHelper) VerifyOTP(keyPrefix, identifier, otp string) error {
	key := fmt.Sprintf("otp:%s:%s", keyPrefix, identifier)
	storedHash, err := r.Client.Get(r.Ctx, key).Result()

	if err == redis.Nil {
		return fmt.Errorf("OTP tidak ditemukan atau sudah kedaluwarsa")
	}
	if err != nil {
		return err
	}

	if storedHash != HashOTP(otp) {
		return fmt.Errorf("OTP tidak cocok")
	}

	// Hapus OTP setelah diverifikasi
	r.Client.Del(r.Ctx, key)

	return nil
}

func (r *RedisHelper) RetryUntilRedisKeyExpired(
	keyPrefix string,
	identifier string,
	retryDelay time.Duration,
	f func() error,
) error {
	key := fmt.Sprintf("otp:%s:%s", keyPrefix, identifier)

	for {
		ttl, err := r.Client.TTL(r.Ctx, key).Result()
		if err != nil {
			return fmt.Errorf("gagal cek TTL Redis: %w", err)
		}

		if ttl <= 0 {
			return fmt.Errorf("OTP sudah kedaluwarsa")
		}

		// Coba kirim ulang
		if err := f(); err != nil {
			log.Printf("Retry kirim OTP gagal, coba lagi dalam %s: %v", retryDelay, err)
			time.Sleep(retryDelay)
			continue
		}

		return nil
	}
}

func (r *RedisHelper) GetOTP(prefix, identifier string) (string, error) {
	key := fmt.Sprintf("otp:%s:%s", prefix, identifier)
	otp, err := r.Client.Get(r.Ctx, key).Result()

	if err == redis.Nil {
		return "", fmt.Errorf("OTP tidak ditemukan atau sudah kedaluwarsa")
	} else if err != nil {
		return "", fmt.Errorf("gagal mengambil OTP: %v", err)
	}

	return otp, nil
}

func (r *RedisHelper) DeleteOTP(prefix, identifier string) error {
	key := fmt.Sprintf("otp:%s:%s", prefix, identifier)
	return r.Client.Del(r.Ctx, key).Err()
}
