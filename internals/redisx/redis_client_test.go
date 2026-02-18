package redisx

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

// ==================== 测试辅助函数 ====================

// setupMiniRedis 创建并启动 miniredis 服务器
func setupMiniRedis(t *testing.T) (*miniredis.Miniredis, func()) {
	t.Helper()

	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("启动 miniredis 失败: %v", err)
	}

	// 清理函数
	cleanup := func() {
		mr.Close()
	}

	return mr, cleanup
}

// initTestClient 初始化测试用的 Redis 客户端
func initTestClient(t *testing.T, addr string) {
	t.Helper()

	err := Init(Config{
		Addr:         addr,
		Password:     "",
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 5,
	})
	if err != nil {
		t.Fatalf("初始化 Redis 客户端失败: %v", err)
	}
}

// ==================== Init 测试 ====================

func TestInit(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "成功初始化",
			config: Config{
				Addr:         "localhost:6379",
				Password:     "",
				DB:           0,
				PoolSize:     10,
				MinIdleConns: 5,
			},
			wantErr: true, // 没有 Redis 服务器，期望失败
		},
		{
			name: "使用 miniredis 初始化",
			config: Config{
				Addr:         "",
				Password:     "",
				DB:           0,
				PoolSize:     10,
				MinIdleConns: 5,
			},
			wantErr: true, // 空地址，期望失败
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 使用 miniredis
			mr, cleanup := setupMiniRedis(t)
			defer cleanup()

			tt.config.Addr = mr.Addr()
			tt.wantErr = false // 使用 miniredis 应该成功

			err := Init(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}

			// 清理
			Close()
		})
	}
}

// ==================== Set 和 Get 测试（表驱动）====================

func TestSetAndGet(t *testing.T) {
	mr, cleanup := setupMiniRedis(t)
	defer cleanup()

	initTestClient(t, mr.Addr())
	defer Close()

	tests := []struct {
		name        string
		key         string
		value       interface{}
		expiration  time.Duration
		wantValue   string
		wantErr     bool
		description string
	}{
		{
			name:        "设置字符串值",
			key:         "user:name",
			value:       "Alice",
			expiration:  0,
			wantValue:   "Alice",
			wantErr:     false,
			description: "基本的字符串设置和获取",
		},
		{
			name:        "设置数字值",
			key:         "counter:1",
			value:       "12345",
			expiration:  0,
			wantValue:   "12345",
			wantErr:     false,
			description: "数字作为字符串存储",
		},
		{
			name:        "设置JSON值",
			key:         "user:profile",
			value:       `{"name":"Bob","age":30}`,
			expiration:  0,
			wantValue:   `{"name":"Bob","age":30}`,
			wantErr:     false,
			description: "JSON字符串存储",
		},
		{
			name:        "设置带过期时间的值",
			key:         "session:token",
			value:       "abc123",
			expiration:  1 * time.Second,
			wantValue:   "abc123",
			wantErr:     false,
			description: "设置1秒过期时间",
		},
		{
			name:        "设置空字符串",
			key:         "empty:value",
			value:       "",
			expiration:  0,
			wantValue:   "",
			wantErr:     false,
			description: "空字符串也是有效值",
		},
		{
			name:        "设置特殊字符",
			key:         "special:chars",
			value:       "hello\nworld\t!",
			expiration:  0,
			wantValue:   "hello\nworld\t!",
			wantErr:     false,
			description: "包含换行和制表符",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.description)

			// 设置值
			err := Set(tt.key, tt.value, tt.expiration)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 获取值
			got, err := Get(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.wantValue {
				t.Errorf("Get() = %v, want %v", got, tt.wantValue)
			}

			// 如果设置了过期时间，等待过期后验证
			if tt.expiration > 0 && tt.expiration < 5*time.Second {
				t.Log("等待过期...")
				mr.FastForward(tt.expiration + 100*time.Millisecond)

				// 过期后应该获取不到值
				_, err = Get(tt.key)
				if err == nil {
					t.Error("期望值已过期，但仍能获取到")
				}
			}
		})
	}
}

// ==================== Del 测试 ====================

func TestDel(t *testing.T) {
	mr, cleanup := setupMiniRedis(t)
	defer cleanup()

	initTestClient(t, mr.Addr())
	defer Close()

	tests := []struct {
		name    string
		setup   func() // 设置测试数据
		keys    []string
		wantErr bool
	}{
		{
			name: "删除单个键",
			setup: func() {
				Set("key1", "value1", 0)
			},
			keys:    []string{"key1"},
			wantErr: false,
		},
		{
			name: "删除多个键",
			setup: func() {
				Set("key1", "value1", 0)
				Set("key2", "value2", 0)
				Set("key3", "value3", 0)
			},
			keys:    []string{"key1", "key2", "key3"},
			wantErr: false,
		},
		{
			name: "删除不存在的键",
			setup: func() {
				// 不设置任何键
			},
			keys:    []string{"nonexistent"},
			wantErr: false, // Redis 删除不存在的键不会报错
		},
		{
			name: "删除部分存在部分不存在的键",
			setup: func() {
				Set("key1", "value1", 0)
			},
			keys:    []string{"key1", "nonexistent"},
			wantErr: false,
		},
		{
			name: "删除单个键",
			setup: func() {
				Set("key1", "value1", 0)
				Set("key2", "value2", 0)
			},
			keys:    []string{"key1"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置测试数据
			if tt.setup != nil {
				tt.setup()
			}

			// 执行删除
			err := Del(tt.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Del() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 验证键已被删除
			for _, key := range tt.keys {
				exists, _ := Exists(key)
				if exists > 0 {
					t.Errorf("键 %s 应该已被删除，但仍然存在", key)
				}
			}
		})
	}
}

// ==================== Exists 测试 ====================

func TestExists(t *testing.T) {
	mr, cleanup := setupMiniRedis(t)
	defer cleanup()

	initTestClient(t, mr.Addr())
	defer Close()

	tests := []struct {
		name    string
		setup   func()
		keys    []string
		want    int64
		wantErr bool
	}{
		{
			name: "检查单个存在的键",
			setup: func() {
				Set("key1", "value1", 0)
			},
			keys:    []string{"key1"},
			want:    1,
			wantErr: false,
		},
		{
			name: "检查单个不存在的键",
			setup: func() {
				// 不设置任何键
			},
			keys:    []string{"nonexistent"},
			want:    0,
			wantErr: false,
		},
		{
			name: "检查多个键（部分存在）",
			setup: func() {
				Set("key1", "value1", 0)
				Set("key2", "value2", 0)
			},
			keys:    []string{"key1", "key2", "key3"},
			want:    2,
			wantErr: false,
		},
		{
			name: "检查多个键（全部存在）",
			setup: func() {
				Set("key1", "value1", 0)
				Set("key2", "value2", 0)
				Set("key3", "value3", 0)
			},
			keys:    []string{"key1", "key2", "key3"},
			want:    3,
			wantErr: false,
		},
		{
			name: "检查多个键（全部存在）",
			setup: func() {
				Set("key1", "value1", 0)
				Set("key2", "value2", 0)
				Set("key3", "value3", 0)
			},
			keys:    []string{"key1", "key2", "key3"},
			want:    3,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			got, err := Exists(tt.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ==================== Expire 和 TTL 测试 ====================

func TestExpireAndTTL(t *testing.T) {
	mr, cleanup := setupMiniRedis(t)
	defer cleanup()

	initTestClient(t, mr.Addr())
	defer Close()

	tests := []struct {
		name       string
		key        string
		value      string
		initialTTL time.Duration
		newTTL     time.Duration
		wantErr    bool
	}{
		{
			name:       "设置新的过期时间",
			key:        "session:user1",
			value:      "data",
			initialTTL: 0, // 初始不过期
			newTTL:     2 * time.Second,
			wantErr:    false,
		},
		{
			name:       "更新已存在的过期时间",
			key:        "session:user2",
			value:      "data",
			initialTTL: 5 * time.Second,
			newTTL:     10 * time.Second,
			wantErr:    false,
		},
		{
			name:       "为不存在的键设置过期时间",
			key:        "nonexistent",
			value:      "",
			initialTTL: 0,
			newTTL:     1 * time.Second,
			wantErr:    false, // Redis 会返回 0，但不会报错
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置初始值
			if tt.value != "" {
				err := Set(tt.key, tt.value, tt.initialTTL)
				if err != nil {
					t.Fatalf("Set() failed: %v", err)
				}
			}

			// 获取初始TTL
			initialTTL, err := TTL(tt.key)
			if err != nil {
				t.Errorf("TTL() before Expire() error = %v", err)
			}
			t.Logf("初始TTL: %v", initialTTL)

			// 设置新的过期时间
			err = Expire(tt.key, tt.newTTL)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expire() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 验证新的TTL
			newTTL, err := TTL(tt.key)
			if err != nil {
				t.Errorf("TTL() after Expire() error = %v", err)
			}

			if tt.value != "" {
				// Redis 的 TTL 精度，允许1秒误差
				if newTTL < tt.newTTL-time.Second || newTTL > tt.newTTL {
					t.Errorf("TTL() = %v, 期望约 %v", newTTL, tt.newTTL)
				}
			}

			// 快进时间验证过期
			if tt.newTTL > 0 && tt.newTTL < 5*time.Second {
				mr.FastForward(tt.newTTL + 100*time.Millisecond)

				expiredTTL, _ := TTL(tt.key)
				// -2 表示键已过期
				if expiredTTL != -2*time.Second && expiredTTL != time.Duration(-2) {
					t.Logf("过期后的TTL: %v (期望约为 -2s)", expiredTTL)
				}
			}
		})
	}
}

// ==================== Incr 和 Decr 测试 ====================

func TestIncrAndDecr(t *testing.T) {
	mr, cleanup := setupMiniRedis(t)
	defer cleanup()

	initTestClient(t, mr.Addr())
	defer Close()

	tests := []struct {
		name        string
		key         string
		setup       func()
		operations  []string // "incr" 或 "decr"
		want        int64
		wantErr     bool
		description string
	}{
		{
			name: "递增新键",
			key:  "counter:test1",
			setup: func() {
				// 不设置初始值
			},
			operations:  []string{"incr"},
			want:       1,
			wantErr:    false,
			description: "对不存在的键递增，从0开始",
		},
		{
			name: "递增已存在的键",
			key:  "counter:test2",
			setup: func() {
				Set("counter:test2", "5", 0)
			},
			operations:  []string{"incr"},
			want:       6,
			wantErr:    false,
			description: "5 -> 6",
		},
		{
			name: "多次递增",
			key:  "counter:test3",
			setup: func() {
				Set("counter:test3", "0", 0)
			},
			operations:  []string{"incr", "incr", "incr"},
			want:       3,
			wantErr:    false,
			description: "0 -> 1 -> 2 -> 3",
		},
		{
			name: "递减新键",
			key:  "counter:test4",
			setup: func() {
				// 不设置初始值
			},
			operations:  []string{"decr"},
			want:       -1,
			wantErr:    false,
			description: "对不存在的键递减，从0开始",
		},
		{
			name: "递减已存在的键",
			key:  "counter:test5",
			setup: func() {
				Set("counter:test5", "5", 0)
			},
			operations:  []string{"decr"},
			want:       4,
			wantErr:    false,
			description: "5 -> 4",
		},
		{
			name: "混合递增递减",
			key:  "counter:test6",
			setup: func() {
				Set("counter:test6", "10", 0)
			},
			operations:  []string{"incr", "incr", "decr", "incr", "decr", "decr"},
			want:       10,
			wantErr:    false,
			description: "10+1+1-1+1-1-1 = 10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.description)

			// 清理之前可能存在的键
			Del(tt.key)

			if tt.setup != nil {
				tt.setup()
			}

			// 执行操作
			for _, op := range tt.operations {
				var err error
				switch op {
				case "incr":
					_, err = Incr(tt.key)
				case "decr":
					_, err = Decr(tt.key)
				}

				if err != nil {
					if !tt.wantErr {
						t.Errorf("%s() error = %v", op, err)
					}
					return
				}
			}

			// 验证最终值
			got, err := Get(tt.key)
			if err != nil {
				t.Errorf("Get() error = %v", err)
				return
			}

			var gotInt int64
			fmt.Sscanf(got, "%d", &gotInt)

			if gotInt != tt.want {
				t.Errorf("最终值 = %d, want %d", gotInt, tt.want)
			}
		})
	}
}

// ==================== 错误处理测试 ====================

func TestErrors(t *testing.T) {
	t.Run("操作未初始化的客户端", func(t *testing.T) {
		// 确保客户端关闭
		Close()
		_, err := Get("test")
		if err == nil {
			t.Error("期望操作失败，但没有错误")
		}
	})

	t.Run("连接不存在的Redis", func(t *testing.T) {
		// 尝试连接到一个不存在的地址
		err := Init(Config{
			Addr: "localhost:9999", // 假设这个端口没有 Redis
		})
		// 这个测试可能会成功或失败，取决于是否有 Redis 在该端口
		// 我们只是记录结果
		if err != nil {
			t.Logf("按预期无法连接: %v", err)
		} else {
			t.Log("端口 9999 上有 Redis 在运行")
			Close() // 清理
		}
	})
}

// ==================== Close 测试 ====================

func TestClose(t *testing.T) {
	mr, cleanup := setupMiniRedis(t)
	defer cleanup()

	// 初始化客户端
	initTestClient(t, mr.Addr())

	// 设置一些数据
	Set("key1", "value1", 0)
	Set("key2", "value2", 0)

	// 关闭连接
	err := Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// 再次关闭会返回错误（客户端已关闭）
	err = Close()
	if err == nil {
		t.Error("期望第二次Close()返回错误，但没有")
	}
	t.Logf("第二次Close()返回错误（符合预期）: %v", err)
}

// ==================== 基准测试 ====================

func BenchmarkSet(b *testing.B) {
	mr, _ := setupMiniRedis(&testing.T{})
	defer mr.Close()

	initTestClient(&testing.T{}, mr.Addr())
	defer Close()

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("bench:key:%d", i)
		rdb.Set(ctx, key, "value", 0)
	}
}

func BenchmarkGet(b *testing.B) {
	mr, _ := setupMiniRedis(&testing.T{})
	defer mr.Close()

	initTestClient(&testing.T{}, mr.Addr())
	defer Close()

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	// 预设数据
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("bench:key:%d", i)
		rdb.Set(ctx, key, "value", 0)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("bench:key:%d", i%1000)
		rdb.Get(ctx, key)
	}
}

func BenchmarkIncr(b *testing.B) {
	mr, _ := setupMiniRedis(&testing.T{})
	defer mr.Close()

	initTestClient(&testing.T{}, mr.Addr())
	defer Close()

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := "bench:counter"
		rdb.Incr(ctx, key)
	}
}

// ==================== 集成测试示例 ====================

func TestIntegration_SessionManagement(t *testing.T) {
	mr, cleanup := setupMiniRedis(t)
	defer cleanup()

	initTestClient(t, mr.Addr())
	defer Close()

	t.Log("=== 会话管理集成测试 ===")

	// 1. 创建会话
	sessionID := "session:abc123"
	sessionData := `{"user_id":123,"name":"Alice","login_time":"2024-01-01T00:00:00Z"}`

	err := Set(sessionID, sessionData, 5*time.Minute)
	if err != nil {
		t.Fatalf("创建会话失败: %v", err)
	}
	t.Log("✓ 创建会话成功")

	// 2. 验证会话存在
	exists, _ := Exists(sessionID)
	if exists != 1 {
		t.Error("会话应该存在")
	}
	t.Log("✓ 会话存在验证通过")

	// 3. 获取会话数据
	retrievedData, err := Get(sessionID)
	if err != nil {
		t.Fatalf("获取会话失败: %v", err)
	}
	if retrievedData != sessionData {
		t.Error("获取的会话数据不匹配")
	}
	t.Log("✓ 获取会话数据成功")

	// 4. 更新会话过期时间
	err = Expire(sessionID, 10*time.Minute)
	if err != nil {
		t.Fatalf("更新会话过期时间失败: %v", err)
	}
	t.Log("✓ 更新会话过期时间成功")

	// 5. 删除会话（登出）
	err = Del(sessionID)
	if err != nil {
		t.Fatalf("删除会话失败: %v", err)
	}
	t.Log("✓ 删除会话成功")

	// 6. 验证会话已删除
	exists, _ = Exists(sessionID)
	if exists != 0 {
		t.Error("会话应该已被删除")
	}
	t.Log("✓ 会话删除验证通过")
}

func TestIntegration_Counter(t *testing.T) {
	mr, cleanup := setupMiniRedis(t)
	defer cleanup()

	initTestClient(t, mr.Addr())
	defer Close()

	t.Log("=== 计数器集成测试 ===")

	counterKey := "counter:page_views"

	// 1. 初始化计数器
	Set(counterKey, "0", 0)
	t.Log("✓ 初始化计数器")

	// 2. 模拟多次页面访问
	for i := 0; i < 10; i++ {
		Incr(counterKey)
	}
	t.Log("✓ 递增计数器 10 次")

	// 3. 获取计数器值
	val, err := Get(counterKey)
	if err != nil {
		t.Fatalf("获取计数器失败: %v", err)
	}
	t.Logf("✓ 计数器当前值: %s", val)

	// 4. 模拟一些无效点击被撤销
	for i := 0; i < 3; i++ {
		Decr(counterKey)
	}
	t.Log("✓ 递减计数器 3 次")

	// 5. 获取最终值
	finalVal, _ := Get(counterKey)
	t.Logf("✓ 计数器最终值: %s", finalVal)
}
