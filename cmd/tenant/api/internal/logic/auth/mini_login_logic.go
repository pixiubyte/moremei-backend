package auth

import (
	"context"
	"errors"
	"moremei/ai-saas/cmd/tenant/api/internal/consts"
	"moremei/ai-saas/cmd/tenant/api/internal/errs"
	"moremei/ai-saas/pkg/gormc"
	"moremei/ai-saas/pkg/jwt"
	"time"

	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MiniLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMiniLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MiniLoginLogic {
	return &MiniLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MiniLoginLogic) MiniLogin(req *types.MiniLoginRequest) (resp *types.LoginResponse, err error) {
	miniClient, err := l.svcCtx.SdkClients.GetMiniClient(l.ctx)
	if err != nil {
		return nil, err
	}
	respOpenid, err := miniClient.GetOpenid(req.JsCode)
	if err != nil {
		l.Logger.Errorf("获取openid失败err:%v", err)
		return nil, err
	}
	miniUser, err := l.svcCtx.MiniUserModel.FindOneByOpenid(l.ctx, respOpenid.Openid)
	if err != nil {
		if !errors.Is(err, gormc.ErrNotFound) {
			return nil, err
		}
		return nil, errs.UserNotRegisterError
	} else if respOpenid.Unionid != "" && respOpenid.Unionid != miniUser.Unionid {
		miniUser.Unionid = respOpenid.Unionid
		err = l.svcCtx.MiniUserModel.Update(l.ctx, miniUser, nil)
		if err != nil {
			l.Errorf("update miniuser err:%v", err)
		}
	}

	user, err := l.svcCtx.UserModel.FindOne(l.ctx, miniUser.UserId)
	if err != nil {
		return nil, err
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

	return &types.LoginResponse{
		AccessToken: token,
		ExpiresIn:   expiresTime,
		UserId:      user.Uuid,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
	}, nil
}
