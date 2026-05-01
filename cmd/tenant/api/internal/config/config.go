package config

import (
	"moremei/ai-saas/pkg/captcha"
	"moremei/ai-saas/pkg/cos"
	"moremei/ai-saas/pkg/gormc"
	"moremei/ai-saas/pkg/uuid"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	Auth      Auth
	Snowflake uuid.SnowflakeConfig
	MySQL     gormc.Config
	Redis     redis.RedisConf
	Cos       cos.Config

	TestAccounts []*TestAccount `json:",optional"`

	Sms struct {
		ErrorLimit     int    `json:",default=5"`
		DevCode        string `json:",default=888888"`
		EnabledCaptcha bool   `json:",default=false"`
	}

	Captcha captcha.Config `json:",optional"`
}
