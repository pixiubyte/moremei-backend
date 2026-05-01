package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	"moremei/ai-saas/cmd/tenant/api/internal/utils"
	"moremei/ai-saas/internal/model/tenant"
	"moremei/ai-saas/pkg/gormc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserAppProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserAppProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserAppProfileLogic {
	return &GetUserAppProfileLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetUserAppProfileLogic) GetUserAppProfile(req *types.UserAppProfileRequest) (*types.UserAppProfileResponse, error) {
	authUser, err := utils.GetAuthUserCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, authUser.UserId)
	if err != nil {
		return nil, err
	}
	if user.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}

	profile, err := l.svcCtx.AppUserProfileModel.FindOneByAppUuidKey(l.ctx, req.AppName, user.Uuid, req.ProfileName)
	if err != nil {
		if !errors.Is(err, gormc.ErrNotFound) {
			return nil, err
		}
		resp := &types.UserAppProfileResponse{ProfileName: req.ProfileName, Profile: struct{}{}}
		sc, err := l.svcCtx.SystemConfigModel.FindByKey(l.ctx, fmt.Sprintf("app_%s_%s_default", req.AppName, req.ProfileName))
		if err != nil {
			return resp, nil
		}
		appuser := tenant.AppUserProfile{App: req.AppName, Uuid: user.Uuid, Key: req.ProfileName, Value: sc.Value}
		if err := l.svcCtx.AppUserProfileModel.Update(l.ctx, &appuser, nil); err != nil {
			l.Logger.Errorf("update app user profile error: %v", err)
			return resp, nil
		}
		profile = &appuser
	}

	resp := &types.UserAppProfileResponse{ProfileName: req.ProfileName, Profile: map[string]interface{}{}}
	if err := json.Unmarshal([]byte(profile.Value), &resp.Profile); err != nil {
		return nil, err
	}
	return resp, nil
}
