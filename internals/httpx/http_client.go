package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// Client HTTP 客户端封装
type Client struct {
	client      *http.Client
	baseURL     string
	headers     map[string]string
	timeout     time.Duration
	retryTimes  int
	retryDelay  time.Duration
}

// Config 客户端配置
type Config struct {
	BaseURL    string
	Timeout    time.Duration
	MaxRetries int
	RetryDelay time.Duration
	Headers    map[string]string
}

// NewClient 创建新的 HTTP 客户端
func NewClient(config Config) *Client {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.RetryDelay == 0 {
		config.RetryDelay = 100 * time.Millisecond
	}

	return &Client{
		client: &http.Client{
			Timeout: config.Timeout,
		},
		baseURL:    config.BaseURL,
		headers:    config.Headers,
		timeout:    config.Timeout,
		retryTimes: config.MaxRetries,
		retryDelay: config.RetryDelay,
	}
}

// Request 通用请求方法
func (c *Client) Request(ctx context.Context, method, path string, body io.Reader, headers map[string]string) (*http.Response, error) {
	url := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	// 设置默认 header
	if c.headers != nil {
		for k, v := range c.headers {
			req.Header.Set(k, v)
		}
	}

	// 设置请求特定的 header
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// 重试逻辑
	var resp *http.Response
	var lastErr error
	for i := 0; i <= c.retryTimes; i++ {
		resp, err = c.client.Do(req)
		if err == nil {
			break
		}
		lastErr = err
		if i < c.retryTimes {
			time.Sleep(c.retryDelay)
		}
	}

	if err != nil {
		return nil, lastErr
	}

	return resp, nil
}

// Get GET 请求
func (c *Client) Get(ctx context.Context, path string, headers map[string]string) (*http.Response, error) {
	return c.Request(ctx, http.MethodGet, path, nil, headers)
}

// Post JSON POST 请求
func (c *Client) Post(ctx context.Context, path string, data interface{}, headers map[string]string) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/json"

	return c.Request(ctx, http.MethodPost, path, bytes.NewReader(jsonData), headers)
}

// Put JSON PUT 请求
func (c *Client) Put(ctx context.Context, path string, data interface{}, headers map[string]string) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/json"

	return c.Request(ctx, http.MethodPut, path, bytes.NewReader(jsonData), headers)
}

// Delete DELETE 请求
func (c *Client) Delete(ctx context.Context, path string, headers map[string]string) (*http.Response, error) {
	return c.Request(ctx, http.MethodDelete, path, nil, headers)
}

// ParseResponse 解析响应体到指定结构
func ParseResponse[T any](resp *http.Response) (T, error) {
	var result T
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	return result, err
}

// ParseRawResponse 解析响应体为原始字节
func ParseRawResponse(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
