package redisx

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

// Config Redis 连接配置
type Config struct {
	Addr         string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
}

// Init 初始化 Redis 客户端
func Init(config Config) error {
	rdb = redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
	})

	return rdb.Ping(ctx).Err()
}

// Set 设置键值对
func Set(key string, value interface{}, expiration time.Duration) error {
	return rdb.Set(ctx, key, value, expiration).Err()
}

// Get 获取键值
func Get(key string) (string, error) {
	return rdb.Get(ctx, key).Result()
}

// Del 删除键
func Del(keys ...string) error {
	return rdb.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func Exists(keys ...string) (int64, error) {
	return rdb.Exists(ctx, keys...).Result()
}

// Expire 设置过期时间
func Expire(key string, expiration time.Duration) error {
	return rdb.Expire(ctx, key, expiration).Err()
}

// TTL 获取剩余过期时间
func TTL(key string) (time.Duration, error) {
	return rdb.TTL(ctx, key).Result()
}

// Incr 自增
func Incr(key string) (int64, error) {
	return rdb.Incr(ctx, key).Result()
}

// Decr 自减
func Decr(key string) (int64, error) {
	return rdb.Decr(ctx, key).Result()
}

// Close 关闭连接
func Close() error {
	if rdb != nil {
		return rdb.Close()
	}
	return nil
}
