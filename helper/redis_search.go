package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func AddProductToAutocomplete(rdb *redis.Client, businessID, productID int, productName string) error {
	ctx := context.Background()

	keyIndex := fmt.Sprintf("autocomplete:product:%d:index", businessID)
	keyData := fmt.Sprintf("autocomplete:product:%d:data", businessID)
	normalized := strings.ToLower(productName)

	product := map[string]interface{}{
		"id":   productID,
		"name": productName,
	}

	jsonValue, err := json.Marshal(product)
	if err != nil {
		return err
	}

	pipe := rdb.TxPipeline()
	pipe.ZAdd(ctx, keyIndex, redis.Z{Score: 0, Member: normalized})
	pipe.HSet(ctx, keyData, normalized, jsonValue)

	_, err = pipe.Exec(ctx)
	return err
}

func UpdateProductAutocomplete(rdb *redis.Client, businessID int, oldName, newName string, productID int) error {
	ctx := context.Background()

	keyIndex := fmt.Sprintf("autocomplete:product:%d:index", businessID)
	keyData := fmt.Sprintf("autocomplete:product:%d:data", businessID)

	oldNorm := strings.ToLower(oldName)
	newNorm := strings.ToLower(newName)

	if oldNorm == newNorm {
		return nil
	}

	pipe := rdb.TxPipeline()
	pipe.ZRem(ctx, keyIndex, oldNorm)
	pipe.HDel(ctx, keyData, oldNorm)

	product := map[string]interface{}{
		"id":   productID,
		"name": newName,
	}

	jsonValue, err := json.Marshal(product)
	if err != nil {
		return err
	}

	pipe.ZAdd(ctx, keyIndex, redis.Z{Score: 0, Member: newNorm})
	pipe.HSet(ctx, keyData, newNorm, jsonValue)

	_, err = pipe.Exec(ctx)
	return err
}

func GetProductAutocomplete(rdb *redis.Client, businessID int, prefix string, limit int64) ([]response.ProductResponse, error) {
	ctx := context.Background()

	keyIndex := fmt.Sprintf("autocomplete:product:%d:index", businessID)
	keyData := fmt.Sprintf("autocomplete:product:%d:data", businessID)

	start := "[" + strings.ToLower(prefix)
	end := "[" + strings.ToLower(prefix) + "\xff"

	names, err := rdb.ZRangeByLex(ctx, keyIndex, &redis.ZRangeBy{
		Min:    start,
		Max:    end,
		Offset: 0,
		Count:  limit,
	}).Result()
	if err != nil {
		return nil, err
	}

	if len(names) == 0 {
		return []response.ProductResponse{}, nil
	}

	jsonList, err := rdb.HMGet(ctx, keyData, names...).Result()
	if err != nil {
		return nil, err
	}

	var results []response.ProductResponse
	for _, raw := range jsonList {
		if raw == nil {
			continue
		}
		var product response.ProductResponse
		if err := json.Unmarshal([]byte(raw.(string)), &product); err == nil {
			results = append(results, product)
		}
		if int64(len(results)) >= limit {
			break
		}
	}

	return results, nil
}

func DeleteProductFromAutocomplete(rdb *redis.Client, businessID int, productName string) error {
	keyIndex := fmt.Sprintf("autocomplete:product:%d:index", businessID)
	normalized := strings.ToLower(productName)

	if err := rdb.ZRem(ctx, keyIndex, normalized).Err(); err != nil {
		return err
	}

	return nil
}

func SetJSONToRedis(ctx context.Context, rdb *redis.Client, key string, data interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, key, jsonData, ttl).Err()
}

func GetJSONFromRedis(ctx context.Context, rdb *redis.Client, key string, dest interface{}) error {
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func DeleteKeysByPattern(ctx context.Context, rdb *redis.Client, pattern string) error {
	iter := rdb.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := rdb.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}
