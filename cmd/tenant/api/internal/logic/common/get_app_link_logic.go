package common

import (
	"context"
	"fmt"
	"time"

	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	miniTypes "moremei/ai-saas/internal/miniprogram/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppLinkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAppLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppLinkLogic {
	return &GetAppLinkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAppLinkLogic) GetAppLink(req *types.GetAppLinkRequest) (resp *types.GetAppLinkResponse, err error) {
	// 使用 redis 缓存，避免重复请求
	cacheKey := fmt.Sprintf("applink:%s:%s:%d:%d:%d:%s:%s", req.AppType, req.Path, req.ExpireType, req.ExpireTime, req.ExpireInterval, req.EnvVersion, req.Query)

	appLink, err := l.svcCtx.RedisClient.Get(cacheKey)
	if err == nil && appLink != "" {
		return &types.GetAppLinkResponse{
			AppLink: appLink,
		}, nil
	}

	miniClient, err := l.svcCtx.SdkClients.GetMiniClient(l.ctx)
	if err != nil {
		return nil, err
	}

	appLinkReq := &miniTypes.GetUrlLinkRequest{
		Path:           req.Path,
		Query:          req.Query,
		ExpireType:     req.ExpireType,
		ExpireTime:     req.ExpireTime,
		ExpireInterval: req.ExpireInterval,
		EnvVersion:     req.EnvVersion,
	}

	appLink, err = miniClient.GetUrlLink(appLinkReq)
	if err != nil {
		return nil, err
	}

	// 缓存
	// 根据传入参数计算时间
	var cacheExpire time.Duration
	if req.ExpireType == 0 {
		// 使用时效时间 expire_time
		if req.ExpireTime > 0 {
			cacheExpire = time.Until(time.Unix(req.ExpireTime, 0))
		} else {
			cacheExpire = time.Hour * 24 * 30 // 默认缓存 30 天
		}
	} else {
		// 使用 expire_interval, 单位为天
		if req.ExpireInterval > 0 && req.ExpireInterval <= 30 {
			cacheExpire = time.Duration(req.ExpireInterval) * 24 * time.Hour
		} else {
			cacheExpire = time.Hour * 24 * 30 // 默认缓存 30 天
		}
	}

	// 转换为 int 类型
	cacheExpireInt := int(cacheExpire.Seconds()) - 60 // 留点安全时间
	if cacheExpireInt < 60 {
		cacheExpireInt = 60 // 最少缓存 1 分钟
	}

	err = l.svcCtx.RedisClient.SetexCtx(l.ctx, cacheKey, appLink, cacheExpireInt)
	if err != nil {
		logx.Errorf("GetAppLink set redis err: %v", err)
	}

	return &types.GetAppLinkResponse{
		AppLink: appLink,
	}, nil
}
