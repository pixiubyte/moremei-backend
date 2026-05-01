package svc

import (
	"context"
	"fmt"

	"moremei/ai-saas/cmd/tenant/api/internal/config"
	"moremei/ai-saas/internal/miniprogram"
	"moremei/ai-saas/internal/model/tenant"
	"moremei/ai-saas/internal/mp"
	"moremei/ai-saas/pkg/captcha"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

type SdkClients struct {
	conf              *config.Config
	systemConfigModel tenant.SystemConfigModel
	redisClient       *redis.Redis

	CaptchaClient *captcha.Client
}

func NewSdkClients(c config.Config, scm tenant.SystemConfigModel, redis *redis.Redis) *SdkClients {
	return &SdkClients{
		conf: &c,

		systemConfigModel: scm,
		redisClient:       redis,

		CaptchaClient: captcha.NewClient(c.Captcha, redis),
	}
}

// 获取小程序客户端
func (c *SdkClients) GetMiniClient(ctx context.Context) (miniprogram.MiniProgram, error) {
	// 从数据库中获取配置
	conf, err := tenant.GetConfigByKey[miniprogram.Config](ctx, c.systemConfigModel, tenant.SystemConfigKeyMiniProgramConf)
	if err != nil {
		return nil, fmt.Errorf("获取小程序配置失败: %v", err)
	}
	// 初始化小程序客户端
	return miniprogram.NewMiniProgram(*conf, c.redisClient, c.conf.Mode), nil
}

// 获取公众号客户端
func (c *SdkClients) GetWxOfficalClient() (mp.MpHandler, error) {
	// 从数据库中获取配置
	conf, err := tenant.GetConfigByKey[mp.Config](context.Background(), c.systemConfigModel, tenant.SystemConfigKeyWxOfficialConf)
	if err != nil {
		return nil, fmt.Errorf("获取公众号配置失败: %v", err)
	}
	// 初始化公众号客户端
	return mp.NewMp(*conf, c.redisClient), nil
}
