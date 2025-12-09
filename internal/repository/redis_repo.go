package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

)

type RedisRepo struct {
	client *redis.Client
}

// NewCacheRepository membuat instance RedisRepo
// Fungsi ini akan dipanggil oleh Wire di cmd/api/wire.go
func NewCacheRepository(client *redis.Client) *RedisRepo {
	return &RedisRepo{client: client}
}

// Get mengambil value berdasarkan key
func (r *RedisRepo) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Set menyimpan value dengan durasi (TTL) tertentu
func (r *RedisRepo) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

// Del menghapus satu key
func (r *RedisRepo) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// DeletePrefix menghapus semua key yang diawali dengan prefix tertentu
// Berguna untuk invalidasi cache, misal: r.DeletePrefix(ctx, "comments:")
func (r *RedisRepo) DeletePrefix(ctx context.Context, prefix string) error {
	iter := r.client.Scan(ctx, 0, prefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		r.client.Del(ctx, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return nil
}