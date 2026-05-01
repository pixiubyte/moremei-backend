package xftech

import (
	"moremei/ai-saas/internal/xftech/cache"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Client struct {
	cfg *Config
	ak  *DefaultAccessToken
}

func MustNewClient(cfg *Config, redisConf redis.RedisConf) *Client {
	ak, err := NewDefaultAccessToken(cfg, cache.NewRedis(redisConf))
	if err != nil {
		logx.Error(err)
	}

	return &Client{
		cfg: cfg,
		ak:  ak,
	}
}
