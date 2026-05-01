package common

import (
	"context"
	"encoding/json"

	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	"moremei/ai-saas/internal/model/tenant"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppVersionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAppVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppVersionLogic {
	return &GetAppVersionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAppVersionLogic) GetAppVersion(req *types.AppVersionRequest) (resp *types.AppVersionResponse, err error) {
	systemConfig, err := l.svcCtx.SystemConfigModel.FindByKey(l.ctx, tenant.SystemConfigKeyAppVersionLatest)
	if err != nil {
		return nil, err
	}

	versions := make(map[string]*types.AppVersionResponse)

	err = json.Unmarshal([]byte(systemConfig.Value), &versions)
	if err != nil {
		return nil, err
	}

	for app, vers := range versions {
		if app == req.AppName {
			resp = vers
			break
		}
	}

	return resp, nil
}
