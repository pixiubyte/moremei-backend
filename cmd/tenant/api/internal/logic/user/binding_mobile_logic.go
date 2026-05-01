package user

import (
	"context"
	"errors"

	"moremei/ai-saas/cmd/tenant/api/internal/errs"
	"moremei/ai-saas/cmd/tenant/api/internal/logic/sms"
	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	"moremei/ai-saas/cmd/tenant/api/internal/utils"
	"moremei/ai-saas/internal/model/tenant"
	"moremei/ai-saas/pkg/gormc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
)

type BindingMobileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBindingMobileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindingMobileLogic {
	return &BindingMobileLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *BindingMobileLogic) BindingMobile(req *types.BindingRequest) error {
	authUser, err := utils.GetAuthUserCtx(l.ctx)
	if err != nil {
		return err
	}
	if l.svcCtx.Config.Mode != service.DevMode {
		if err := sms.NewSmsLogic(l.ctx, l.svcCtx).Verify(req.Mobile, req.Code, tenant.SmsBizTypeBind); err != nil {
			return err
		}
	}

	webUserModel := l.svcCtx.WebUserModel
	originWebUser, err := webUserModel.FindOneByMobile(l.ctx, req.Mobile)
	if err != nil {
		if !errors.Is(err, gormc.ErrNotFound) {
			return err
		}
		return webUserModel.Insert(l.ctx, &tenant.WebUser{Mobile: req.Mobile, UserId: authUser.UserId}, nil)
	}
	if originWebUser.UserId == authUser.UserId {
		return nil
	}
	return errs.UserBindExistMobileError
}
