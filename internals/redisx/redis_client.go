package redisx

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client

func main() {
	// 1. 创建客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 地址
		Password: "",               // 密码（无则留空）
		DB:       0,                // 默认数据库
	})

	// 2. 写入
	err := rdb.Set(ctx, "foo", "bar", 0).Err()
	if err != nil {
		panic(err)
	}

	// 3. 读取
	val, err := rdb.Get(ctx, "foo").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("foo:", val)

	// 4. 不存在的 key
	val2, err := rdb.Get(ctx, "missing").Result()
	if err == redis.Nil {
		fmt.Println("missing does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("missing:", val2)
	}
}

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:         "localhost:6379", // Redis 地址
		Password:     "",               // 密码（无则留空）
		DB:           0,                // 默认数据库
		PoolSize:     20,               // 连接池大小
		MinIdleConns: 5,                // 保持的最小空闲连接数
	})
}

func Set(key string, value string) {
	err := rdb.Set(ctx, "foo", "bar", 0).Err()
	if err != nil {
		panic(err)
	}
}
