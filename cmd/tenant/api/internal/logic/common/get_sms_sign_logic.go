package common

import (
	"context"

	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	"moremei/ai-saas/pkg/captcha"
	errorx "moremei/ai-saas/pkg/x/error"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSmsSignLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSmsSignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSmsSignLogic {
	return &GetSmsSignLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSmsSignLogic) GetSmsSign(req *types.GetSmsSignRequest) (resp *types.GetSmsSignResponse, err error) {
	if req.Scene != captcha.SceneMiniProgram.String() && req.Scene != captcha.SceneWeb.String() {
		return nil, errorx.InvalidParamsError
	}

	encryptCaptchaAppId, err := l.svcCtx.SdkClients.CaptchaClient.EncryptCaptchaAppId()
	if err != nil {
		return nil, err
	}

	return &types.GetSmsSignResponse{Signure: encryptCaptchaAppId}, nil
}
