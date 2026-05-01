package common

import (
	"context"

	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	"moremei/ai-saas/internal/model/tenant"

	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConfigLogic {
	return &GetConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetConfigLogic) GetConfig(req *types.GetConfigRequest) (resp *types.GetConfigResponse, err error) {
	// 查询公开的配置键
	keys, err := tenant.GetConfigByKey[[]string](l.ctx, l.svcCtx.SystemConfigModel, tenant.SystemConfigKeyPublicConfigKeys)
	if err != nil {
		l.Logger.Errorf("GetConfigByKey failed, key: %v; err: %v", tenant.SystemConfigKeyPublicConfigKeys, err)
		return nil, nil // 忽略错误，返回空
	}

	// 判断请求的key是否在公开的配置键中
	if !lo.Contains(*keys, req.Configkey) {
		return nil, nil
	}

	// 查询配置
	config, err := tenant.GetConfigByKey[string](l.ctx, l.svcCtx.SystemConfigModel, req.Configkey)
	if err != nil {
		l.Logger.Errorf("GetConfigByKey failed, key: %v; err: %v", req.Configkey, err)
		return nil, nil // 忽略错误，返回空
	}
	resp = &types.GetConfigResponse{
		Content: *config,
	}

	return
}
