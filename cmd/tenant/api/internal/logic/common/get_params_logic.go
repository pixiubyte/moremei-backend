package common

import (
	"context"
	"errors"
	"time"

	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetParamsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetParamsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetParamsLogic {
	return &GetParamsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetParamsLogic) GetParams(req *types.GetParamsRequest) (resp *types.GetParamsResponse, err error) {
	params, err := l.svcCtx.ParamsStorageModel.FindOne(l.ctx, req.ParamId)
	if err != nil {
		return nil, err
	}

	// 检查是否过期
	if params.ExpireAt.Valid && params.ExpireAt.Time.Before(time.Now()) {
		return nil, errors.New("参数已过期")
	}

	// 更新访问次数
	params.AccessCount++
	err = l.svcCtx.ParamsStorageModel.Update(l.ctx, params, nil)
	if err != nil {
		return nil, err
	}

	resp = &types.GetParamsResponse{
		Params:      params.Params,
		CreatedAt:   params.CreatedAt.Unix(),
		Description: params.Description.String,
	}
	if params.ExpireAt.Valid {
		resp.ExpireAt = params.ExpireAt.Time.Unix()
	}
	return
}
