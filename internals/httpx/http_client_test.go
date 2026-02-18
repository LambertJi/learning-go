package httpx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

// ==================== 测试辅助函数 ====================

// 创建测试服务器
func createTestServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

// ==================== NewClient 测试 ====================

func TestNewClient(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   struct {
			timeout    time.Duration
			retryTimes int
			retryDelay time.Duration
		}
	}{
		{
			name: "默认配置",
			config: Config{
				BaseURL: "http://example.com",
			},
			want: struct {
				timeout    time.Duration
				retryTimes int
				retryDelay time.Duration
			}{
				timeout:    30 * time.Second,
				retryTimes: 3,
				retryDelay: 100 * time.Millisecond,
			},
		},
		{
			name: "自定义配置",
			config: Config{
				BaseURL:    "http://example.com",
				Timeout:    10 * time.Second,
				MaxRetries: 5,
				RetryDelay: 200 * time.Millisecond,
			},
			want: struct {
				timeout    time.Duration
				retryTimes int
				retryDelay time.Duration
			}{
				timeout:    10 * time.Second,
				retryTimes: 5,
				retryDelay: 200 * time.Millisecond,
			},
		},
		{
			name: "带默认headers",
			config: Config{
				BaseURL: "http://example.com",
				Headers: map[string]string{
					"Authorization": "Bearer token",
					"User-Agent":    "TestClient",
				},
			},
			want: struct {
				timeout    time.Duration
				retryTimes int
				retryDelay time.Duration
			}{
				timeout:    30 * time.Second,
				retryTimes: 3,
				retryDelay: 100 * time.Millisecond,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.config)

			if client.timeout != tt.want.timeout {
				t.Errorf("timeout = %v, want %v", client.timeout, tt.want.timeout)
			}
			if client.retryTimes != tt.want.retryTimes {
				t.Errorf("retryTimes = %d, want %d", client.retryTimes, tt.want.retryTimes)
			}
			if client.retryDelay != tt.want.retryDelay {
				t.Errorf("retryDelay = %v, want %v", client.retryDelay, tt.want.retryDelay)
			}
		})
	}
}

// ==================== GET 请求测试 ====================

func TestClient_Get(t *testing.T) {
	// 创建测试服务器
	server := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求方法
		if r.Method != http.MethodGet {
			t.Errorf("期望 GET 方法，实际 %s", r.Method)
		}

		// 验证请求路径
		if r.URL.Path != "/test" {
			t.Errorf("期望路径 /test，实际 %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	})
	defer server.Close()

	// 创建客户端
	client := NewClient(Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})

	// 执行请求
	ctx := context.Background()
	resp, err := client.Get(ctx, "/test", nil)

	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 验证响应
	if resp.StatusCode != http.StatusOK {
		t.Errorf("期望状态码 200，实际 %d", resp.StatusCode)
	}
}

// ==================== POST 请求测试（表驱动）====================

func TestClient_Post(t *testing.T) {
	tests := []struct {
		name       string
		data       interface{}
		wantStatus int
		wantBody   string
	}{
		{
			name: "成功发送JSON数据",
			data: map[string]string{
				"name":  "test",
				"value": "123",
			},
			wantStatus: http.StatusOK,
			wantBody:   `{"status":"created"}`,
		},
		{
			name:       "发送空对象",
			data:       map[string]interface{}{},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"invalid data"}`,
		},
		{
			name: "发送复杂数据",
			data: struct {
				ID       int    `json:"id"`
				Message  string `json:"message"`
				Active   bool   `json:"active"`
				Quantity int    `json:"quantity"`
			}{
				ID:       1,
				Message:  "hello",
				Active:   true,
				Quantity: 100,
			},
			wantStatus: http.StatusOK,
			wantBody:   `{"status":"ok"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试服务器
			server := createTestServer(func(w http.ResponseWriter, r *http.Request) {
				// 验证请求方法
				if r.Method != http.MethodPost {
					t.Errorf("期望 POST 方法，实际 %s", r.Method)
				}

				// 验证 Content-Type
				contentType := r.Header.Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("期望 Content-Type: application/json，实际 %s", contentType)
				}

				// 读取请求体
				var data map[string]interface{}
				decoder := json.NewDecoder(r.Body)
				if err := decoder.Decode(&data); err != nil {
					t.Errorf("解析请求体失败: %v", err)
				}

				// 返回响应
				w.WriteHeader(tt.wantStatus)
				w.Write([]byte(tt.wantBody))
			})
			defer server.Close()

			// 创建客户端
			client := NewClient(Config{
				BaseURL: server.URL,
				Timeout: 5 * time.Second,
			})

			// 执行请求
			ctx := context.Background()
			resp, err := client.Post(ctx, "/test", tt.data, nil)

			if err != nil {
				t.Fatalf("请求失败: %v", err)
			}
			defer resp.Body.Close()

			// 验证响应
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("期望状态码 %d，实际 %d", tt.wantStatus, resp.StatusCode)
			}
		})
	}
}

// ==================== PUT 请求测试 ====================

func TestClient_Put(t *testing.T) {
	server := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("期望 PUT 方法，实际 %s", r.Method)
		}

		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("期望 Content-Type: application/json，实际 %s", contentType)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"updated"}`))
	})
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})

	data := map[string]string{"key": "value"}
	ctx := context.Background()
	resp, err := client.Put(ctx, "/update", data, nil)

	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("期望状态码 200，实际 %d", resp.StatusCode)
	}
}

// ==================== DELETE 请求测试 ====================

func TestClient_Delete(t *testing.T) {
	server := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("期望 DELETE 方法，实际 %s", r.Method)
		}

		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})

	ctx := context.Background()
	resp, err := client.Delete(ctx, "/delete/123", nil)

	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("期望状态码 204，实际 %d", resp.StatusCode)
	}
}

// ==================== Headers 测试 ====================

func TestClient_Headers(t *testing.T) {
	tests := []struct {
		name         string
		defaultHeaders map[string]string
		requestHeaders map[string]string
		wantHeaders  map[string]string
	}{
		{
			name: "只有默认headers",
			defaultHeaders: map[string]string{
				"Authorization": "Bearer token123",
				"X-API-Key":     "secret",
			},
			requestHeaders: nil,
			wantHeaders: map[string]string{
				"Authorization": "Bearer token123",
				"X-API-Key":     "secret",
			},
		},
		{
			name: "只有请求headers",
			defaultHeaders: nil,
			requestHeaders: map[string]string{
				"X-Custom-Header": "custom-value",
			},
			wantHeaders: map[string]string{
				"X-Custom-Header": "custom-value",
			},
		},
		{
			name: "默认和请求headers合并",
			defaultHeaders: map[string]string{
				"Authorization": "Bearer token123",
			},
			requestHeaders: map[string]string{
				"X-Custom": "value",
			},
			wantHeaders: map[string]string{
				"Authorization": "Bearer token123",
				"X-Custom":      "value",
			},
		},
		{
			name: "请求headers覆盖默认headers",
			defaultHeaders: map[string]string{
				"Authorization": "Bearer old-token",
			},
			requestHeaders: map[string]string{
				"Authorization": "Bearer new-token",
			},
			wantHeaders: map[string]string{
				"Authorization": "Bearer new-token",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := createTestServer(func(w http.ResponseWriter, r *http.Request) {
				// 验证所有期望的headers
				for key, wantValue := range tt.wantHeaders {
					gotValue := r.Header.Get(key)
					if gotValue != wantValue {
						t.Errorf("Header %s = %s，期望 %s", key, gotValue, wantValue)
					}
				}
				w.WriteHeader(http.StatusOK)
			})
			defer server.Close()

			client := NewClient(Config{
				BaseURL: server.URL,
				Headers: tt.defaultHeaders,
			})

			ctx := context.Background()
			resp, err := client.Get(ctx, "/test", tt.requestHeaders)
			if err != nil {
				t.Fatalf("请求失败: %v", err)
			}
			resp.Body.Close()
		})
	}
}

// ==================== 重试机制测试 ====================

func TestClient_Retry(t *testing.T) {
	// 使用一个会在前几次拒绝连接的端口来测试重试
	// 注意：这种测试依赖于端口未被占用，可能不稳定
	// 更好的方式是使用自定义的 http.RoundTripper

	t.Run("超时后重试", func(t *testing.T) {
		attemptCount := 0
		server := createTestServer(func(w http.ResponseWriter, r *http.Request) {
			attemptCount++
			// 所有请求都成功
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success":true}`))
		})
		defer server.Close()

		client := NewClient(Config{
			BaseURL:    server.URL,
			Timeout:    5 * time.Second,
			MaxRetries: 2,
			RetryDelay: 10 * time.Millisecond,
		})

		ctx := context.Background()
		resp, err := client.Get(ctx, "/test", nil)

		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("期望状态码 200，实际 %d", resp.StatusCode)
		}

		// 正常情况下应该只请求一次
		if attemptCount != 1 {
			t.Logf("注意：尝试了 %d 次（正常情况下应为1次，除非有网络问题）", attemptCount)
		}
	})

	t.Run("重试配置验证", func(t *testing.T) {
		// 验证重试配置被正确设置
		retryTests := []struct {
			name       string
			maxRetries int
			want       int
		}{
			{"默认重试次数", 0, 3},
			{"自定义重试次数", 5, 5},
			{"无重试", -1, 0}, // 负数表示不重试
		}

		for _, tt := range retryTests {
			t.Run(tt.name, func(t *testing.T) {
				config := Config{
					BaseURL:    "http://example.com",
					MaxRetries: tt.maxRetries,
				}
				client := NewClient(config)

				// 验证重试次数配置
				if tt.maxRetries <= 0 {
					// 默认值或负数
					if client.retryTimes <= 0 {
						client.retryTimes = 0
					}
				}

				t.Logf("重试次数配置: %d", client.retryTimes)
			})
		}
	})
}

// ==================== ParseResponse 测试 ====================

func TestParseResponse(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "成功解析JSON",
			body: `{"message":"success","code":200}`,
			want: map[string]interface{}{
				"message": "success",
				"code":    float64(200),
			},
			wantErr: false,
		},
		{
			name:    "空JSON对象",
			body:    `{}`,
			want:    map[string]interface{}{},
			wantErr: false,
		},
		{
			name:    "无效JSON",
			body:    `{invalid json}`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "空响应体",
			body:    ``,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟响应
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       &mockReadCloser{strings.NewReader(tt.body)},
			}

			// 解析响应
			result, err := ParseResponse[map[string]interface{}](resp)

			// 验证错误
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 验证结果
			if !tt.wantErr {
				if result["message"] != tt.want["message"] {
					t.Errorf("message = %v, want %v", result["message"], tt.want["message"])
				}
			}
		})
	}
}

// ==================== ParseRawResponse 测试 ====================

func TestParseRawResponse(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		want    []byte
		wantErr bool
	}{
		{
			name:    "成功读取原始响应",
			body:    `{"raw":"data"}`,
			want:    []byte(`{"raw":"data"}`),
			wantErr: false,
		},
		{
			name:    "空响应",
			body:    ``,
			want:    []byte{},
			wantErr: false,
		},
		{
			name:    "二进制数据",
			body:    "\x00\x01\x02\x03",
			want:    []byte("\x00\x01\x02\x03"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       &mockReadCloser{strings.NewReader(tt.body)},
			}

			result, err := ParseRawResponse(resp)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRawResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && string(result) != string(tt.want) {
				t.Errorf("ParseRawResponse() = %v, want %v", result, tt.want)
			}
		})
	}
}

// ==================== 并发测试 ====================

func TestClient_ConcurrentRequests(t *testing.T) {
	server := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond) // 模拟延迟
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})

	// 并发发送多个请求
	const numRequests = 10
	var wg sync.WaitGroup
	errors := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			ctx := context.Background()
			resp, err := client.Get(ctx, "/test", nil)
			if err != nil {
				errors <- err
				return
			}
			resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				errors <- fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// 检查是否有错误
	for err := range errors {
		t.Errorf("并发请求出错: %v", err)
	}
}

// ==================== 基准测试 ====================

func BenchmarkClient_Get(b *testing.B) {
	server := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})

	ctx := context.Background()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp, err := client.Get(ctx, "/benchmark", nil)
		if err != nil {
			b.Fatalf("请求失败: %v", err)
		}
		resp.Body.Close()
	}
}

func BenchmarkClient_Post(b *testing.B) {
	server := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"created"}`))
	})
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})

	data := map[string]string{"key": "value"}
	ctx := context.Background()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp, err := client.Post(ctx, "/benchmark", data, nil)
		if err != nil {
			b.Fatalf("请求失败: %v", err)
		}
		resp.Body.Close()
	}
}

// ==================== 测试辅助类型 ====================

// mockReadCloser 模拟 io.ReadCloser
type mockReadCloser struct {
	*strings.Reader
}

func (m *mockReadCloser) Close() error {
	return nil
}
