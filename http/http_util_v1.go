package http

import (
	"context"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

// Client 封装的 HTTP 客户端
type Client struct {
	client *resty.Client
}

// Config HTTP 客户端配置
type Config struct {
	BaseURL    string            // 基础 URL
	Timeout    time.Duration     // 请求超时时间，默认 30 秒
	Headers    map[string]string // 默认请求头
	RetryCount int               // 重试次数，默认 0（不重试）
	RetryWait  time.Duration     // 重试等待时间，默认 1 秒
}

// 默认配置常量
const (
	DefaultTimeout    = 30 * time.Second // 默认超时时间
	DefaultRetryWait  = 1 * time.Second  // 默认重试等待时间
	DefaultRetryCount = 0                // 默认重试次数（不重试）
)

// NewClient 创建新的 HTTP 客户端
func NewClient(cfg Config) *Client {
	client := resty.New()

	// 设置默认值
	if cfg.Timeout == 0 {
		cfg.Timeout = DefaultTimeout
	}
	if cfg.RetryWait == 0 {
		cfg.RetryWait = DefaultRetryWait
	}
	if cfg.RetryCount == 0 {
		cfg.RetryCount = DefaultRetryCount
	}

	// 配置客户端
	client.SetTimeout(cfg.Timeout)
	if cfg.BaseURL != "" {
		client.SetBaseURL(cfg.BaseURL)
	}
	if len(cfg.Headers) > 0 {
		client.SetHeaders(cfg.Headers)
	}
	if cfg.RetryCount > 0 {
		client.SetRetryCount(cfg.RetryCount)
		client.SetRetryWaitTime(cfg.RetryWait)
	}

	return &Client{client: client}
}

// GetRestyClient 返回底层的 resty 客户端
func (c *Client) GetRestyClient() *resty.Client {
	return c.client
}

// GetHTTPClient 返回底层的 net/http 客户端
func (c *Client) GetHTTPClient() *http.Client {
	return c.client.GetClient()
}

// Get 发送 GET 请求
func (c *Client) Get(ctx context.Context, url string, queryParams map[string]string, result any) (*resty.Response, error) {
	req := c.client.R().
		SetContext(ctx).
		SetQueryParams(queryParams).
		SetResult(result)
	return req.Get(url)
}

// Post 发送 POST 请求
func (c *Client) Post(ctx context.Context, url string, body, result any) (*resty.Response, error) {
	req := c.client.R().
		SetContext(ctx).
		SetBody(body).
		SetResult(result)
	return req.Post(url)
}
