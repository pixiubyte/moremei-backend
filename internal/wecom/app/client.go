package app

import "moremei/ai-saas/internal/wecom/context"

// Client 接口实例
type Client struct {
	*context.Context
}

// NewClient 初始化实例
func NewClient(ctx *context.Context) *Client {
	return &Client{ctx}
}
