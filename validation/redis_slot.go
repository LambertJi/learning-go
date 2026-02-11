package validation

import (
	"fmt"
	"learning-go/internals/redisx"
)

func TestRedisSlotStrategy() {
	// 初始化 Redis 连接
	err := redisx.Init(redisx.Config{
		Addr:         "localhost:6379",
		Password:     "",
		DB:           0,
		PoolSize:     20,
		MinIdleConns: 5,
	})
	if err != nil {
		fmt.Printf("Redis 初始化失败: %v\n", err)
		return
	}
	defer redisx.Close()

	// 测试基本操作
	testBasicOps()

	// 性能测试
	keyCount := 100000
	fmt.Printf("开始写入 %d 个 key...\n", keyCount)
	for i := 0; i < keyCount; i++ {
		key := fmt.Sprintf("slot:test:%d", i)
		if err := redisx.Set(key, fmt.Sprintf("value-%d", i), 0); err != nil {
			fmt.Printf("写入失败: %v\n", err)
			return
		}
	}
	fmt.Printf("写入完成\n")

	// 验证写入
	count, _ := redisx.Exists("slot:test:0")
	fmt.Printf("验证结果: 存在的 key 数量 = %d\n", count)
}

func testBasicOps() {
	// Set
	if err := redisx.Set("test:key", "test-value", 0); err != nil {
		fmt.Printf("Set 失败: %v\n", err)
		return
	}

	// Get
	val, err := redisx.Get("test:key")
	if err != nil {
		fmt.Printf("Get 失败: %v\n", err)
		return
	}
	fmt.Printf("Get 结果: %s\n", val)

	// Incr
	count, err := redisx.Incr("test:counter")
	if err != nil {
		fmt.Printf("Incr 失败: %v\n", err)
		return
	}
	fmt.Printf("Incr 结果: %d\n", count)
}
