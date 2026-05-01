package auth

import (
	"context"
	"errors"
	"time"

	"moremei/ai-saas/cmd/tenant/api/internal/consts"
	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	"moremei/ai-saas/internal/model/tenant"
	commonConst "moremei/ai-saas/pkg/consts"
	"moremei/ai-saas/pkg/gormc"
	"moremei/ai-saas/pkg/jwt"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenLogic {
	return &GetTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTokenLogic) GetToken(req *types.GetTokenRequest) (resp *types.GetTokenResponse, err error) {
	var user *tenant.User
	webUser, err := l.svcCtx.WebUserModel.FindOneByMobile(l.ctx, req.Mobile)
	if err != nil {
		if !errors.Is(err, gormc.ErrNotFound) {
			return nil, err
		}
		user = &tenant.User{
			Uuid:       l.svcCtx.Snowflake.Generate(),
			Nickname:   req.NickName,
			Avatar:     req.AvatarUrl,
			Gender:     req.Gender,
			Status:     commonConst.FlagTrue,
			RegisterAt: time.Now(),
		}
		webUser = &tenant.WebUser{
			Mobile:  req.Mobile,
			Deleted: 0,
			From:    1,
		}
		err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
			err = l.svcCtx.UserModel.Insert(l.ctx, user, tx)
			if err != nil {
				return err
			}
			webUser.UserId = user.Id
			return l.svcCtx.WebUserModel.Insert(l.ctx, webUser, tx)
		})
		if err != nil {
			return nil, err
		}
	}
	if user == nil {
		user, err = l.svcCtx.UserModel.FindOne(l.ctx, webUser.UserId)
		if err != nil {
			return nil, err
		}
	}
	expiresTime := time.Now().Add(time.Duration(l.svcCtx.Config.Auth.AccessExpire) * time.Second).Unix()
	jc := jwt.NewClient()
	token, err := jc.SetSecret(l.svcCtx.Config.Auth.AccessSecret).
		SetExpire(l.svcCtx.Config.Auth.AccessExpire).
		GenerateToken(map[string]interface{}{
			consts.UserIDCtxKey:   user.Id,
		})
	if err != nil {
		return nil, err
	}
	return &types.GetTokenResponse{
		AccessToken: token,
		ExpiresIn:   expiresTime,
		UserId:      user.Uuid,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
	}, nil
}
