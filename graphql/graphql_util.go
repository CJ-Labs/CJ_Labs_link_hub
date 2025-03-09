package graphql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/shurcooL/graphql"

	"github.com/CJ-Labs/CJ_Labs_link_hub/http" // 使用正确的包路径
)

// Client 封装的 GraphQL 客户端
type Client struct {
	graphqlClient *graphql.Client
	httpClient    *http.Client
	logger        *log.Logger
	endpoint      string
}

// Config GraphQL 客户端配置
type Config struct {
	Endpoint   string
	HTTPConfig http.Config
	Logger     *log.Logger
}

// 重试相关常量
const (
	DefaultRetryCount     = 3               // 默认重试次数
	DefaultRetryWait      = 1 * time.Second // 默认重试等待时间
	DefaultRetryMaxWait   = 5 * time.Second // 默认最大重试等待时间
	DefaultServerErrorMin = 500             // 服务器错误的最小状态码
)

// NewClient 创建新的 GraphQL 客户端
func NewClient(cfg Config) (*Client, error) {
	if cfg.Endpoint == "" {
		return nil, fmt.Errorf("endpoint cannot be empty")
	}

	// 设置默认日志
	if cfg.Logger == nil {
		cfg.Logger = log.Default()
	}

	// 创建 HTTP 客户端
	httpClient := http.NewClient(cfg.HTTPConfig)

	// 创建 GraphQL 客户端
	graphqlClient := graphql.NewClient(cfg.Endpoint, httpClient.GetHTTPClient())

	return &Client{
		graphqlClient: graphqlClient,
		httpClient:    httpClient,
		logger:        cfg.Logger,
		endpoint:      cfg.Endpoint,
	}, nil
}

// Query 执行 GraphQL 查询
func (c *Client) Query(ctx context.Context, query any, variables map[string]interface{}) error {
	return c.execute(ctx, func(ctx context.Context) error {
		return c.graphqlClient.Query(ctx, query, variables)
	})
}

// Mutate 执行 GraphQL 变更
func (c *Client) Mutate(ctx context.Context, mutation any, variables map[string]interface{}) error {
	return c.execute(ctx, func(ctx context.Context) error {
		return c.graphqlClient.Mutate(ctx, mutation, variables)
	})
}

// execute 执行带有重试逻辑的操作
func (c *Client) execute(ctx context.Context, op func(context.Context) error) error {
	restyClient := c.httpClient.GetRestyClient()

	// 配置重试策略
	restyClient.
		SetRetryCount(DefaultRetryCount).
		SetRetryWaitTime(DefaultRetryWait).
		SetRetryMaxWaitTime(DefaultRetryMaxWait).
		AddRetryCondition(func(resp *resty.Response, err error) bool {
			if err != nil {
				c.logger.Printf("retrying due to error: %v", err)
				return true
			}
			if resp != nil && resp.StatusCode() >= DefaultServerErrorMin {
				c.logger.Printf("retrying due to server error: %d", resp.StatusCode())
				return true
			}
			return false
		})

	// 执行操作
	if err := op(ctx); err != nil {
		return fmt.Errorf("graphql operation failed: %w", err)
	}
	return nil
}

// GetRestyClient 获取底层的 resty 客户端
func (c *Client) GetRestyClient() *resty.Client {
	return c.httpClient.GetRestyClient()
}

// GetHTTPClient 获取底层的 HTTP 客户端
func (c *Client) GetHTTPClient() *http.Client {
	return c.httpClient
}
