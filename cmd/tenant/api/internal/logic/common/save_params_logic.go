package common

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	"moremei/ai-saas/cmd/tenant/api/internal/utils"
	"moremei/ai-saas/internal/model/tenant"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveParamsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveParamsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveParamsLogic {
	return &SaveParamsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveParamsLogic) SaveParams(req *types.SaveParamsRequest) (resp *types.SaveParamsResponse, err error) {
	authUser, err := utils.GetAuthUserCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	params := &tenant.ParamsStorage{
		UserId: authUser.UserId,
		Params: req.Params,
	}
	if req.ExpireAt > 0 {
		// 校验时间是否大于当前时间
		currentTime := time.Now().Unix()
		if req.ExpireAt <= currentTime {
			return nil, errors.New("过期时间必须大于当前时间")
		}

		params.ExpireAt = sql.NullTime{
			Time:  time.Unix(req.ExpireAt, 0), // 秒,
			Valid: true,
		}
	}
	if req.Description != "" {
		params.Description = sql.NullString{
			String: req.Description,
			Valid:  true,
		}
	}
	err = l.svcCtx.ParamsStorageModel.Insert(l.ctx, params, nil)
	if err != nil {
		return nil, err
	}
	return &types.SaveParamsResponse{
		ParamId: params.Id,
	}, nil
}
