package common

import (
	"context"
	"fmt"

	"moremei/ai-saas/cmd/tenant/api/internal/errs"
	"moremei/ai-saas/cmd/tenant/api/internal/logic/sms"
	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	"moremei/ai-saas/cmd/tenant/api/internal/utils"
	"moremei/ai-saas/internal/model/tenant"
	"moremei/ai-saas/pkg/captcha"
	errorx "moremei/ai-saas/pkg/x/error"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type SendSmsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendSmsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendSmsLogic {
	return &SendSmsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendSmsLogic) SendSms(req *types.SmsSendRequest) error {
	clientIp := utils.GetRemoteIpCtx(l.ctx)

	if l.svcCtx.Config.Sms.EnabledCaptcha {
		if req.Ticket == "" || req.Scene == "" {
			return errorx.InvalidParamsError
		}
		verifyOptions := []captcha.VerifyOption{
			captcha.WithTicket(req.Ticket),
			captcha.WithUserIp(clientIp),
			captcha.WithTicketCacheKeyPrefix(fmt.Sprintf("%s:", l.svcCtx.Config.Name)),
		}
		if req.Scene == captcha.SceneWeb.String() {
			if req.RandStr == "" {
				return errorx.InvalidParamsError
			}
			verifyOptions = append(verifyOptions, captcha.WithRandStr(req.RandStr))
		}

		verified, err := l.svcCtx.SdkClients.CaptchaClient.VerifyTicket(captcha.Scene(req.Scene), verifyOptions...)
		if err != nil {
			return err
		}

		if !verified {
			return errs.SmsCaptchaError
		}
	}

	switch req.Method {
	case tenant.SmsBizTypeLogin, tenant.SmsBizTypeBind:
		return sms.NewSmsLogic(l.ctx, l.svcCtx).Send(req.Method, req.Mobile, clientIp)
	default:
		return errors.Errorf("不支持该业务类型:%s", req.Method)
	}
}
