package user

import (
	"context"

	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	"moremei/ai-saas/cmd/tenant/api/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetUserHeightLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetUserHeightLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUserHeightLogic {
	return &SetUserHeightLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *SetUserHeightLogic) SetUserHeight(req *types.SetUserHeightRequest) error {
	authUser, err := utils.GetAuthUserCtx(l.ctx)
	if err != nil {
		return err
	}
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, authUser.UserId)
	if err != nil {
		return err
	}
	user.Height = req.Height
	if err := l.svcCtx.UserModel.Update(l.ctx, user, nil); err != nil {
		l.Logger.Errorf("update user height err:%v", err)
		return err
	}
	return nil
}
