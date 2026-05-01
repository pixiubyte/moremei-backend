package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"moremei/ai-saas/cmd/tenant/api/internal/consts"
	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	"moremei/ai-saas/cmd/tenant/api/internal/utils"
	"moremei/ai-saas/internal/model/tenant"
	pkgUtils "moremei/ai-saas/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
)

type SetUserAppProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetUserAppProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUserAppProfileLogic {
	return &SetUserAppProfileLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *SetUserAppProfileLogic) SetUserAppProfile(req *types.SetUserAppProfileRequest) error {
	authUser, err := utils.GetAuthUserCtx(l.ctx)
	if err != nil {
		return err
	}
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, authUser.UserId)
	if err != nil {
		return err
	}
	if user.Status != 1 {
		return errors.New("用户已被禁用")
	}

	pf, err := json.Marshal(req.Profile)
	if err != nil {
		return err
	}
	appuser := tenant.AppUserProfile{App: req.AppName, Uuid: user.Uuid, Key: req.ProfileName, Value: string(pf)}

	userAidLock := fmt.Sprintf("%s:%s:%v:%v", l.svcCtx.Config.Name, req.AppName, consts.UserProfileKeyArchive, authUser.UserId)
	lock := redis.NewRedisLock(l.svcCtx.RedisClient, userAidLock)
	if err := pkgUtils.AcquireAndWaitRedisLock(l.ctx, lock, 5); err != nil {
		return err
	}
	defer lock.ReleaseCtx(l.ctx)

	apy, err := l.svcCtx.AppUserProfileModel.FindOneByAppUuidKey(l.ctx, appuser.App, user.Uuid, appuser.Key)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	} else {
		appuser.Id = apy.Id
		appuser.CreatedAt = apy.CreatedAt
		appuser.UpdatedAt = apy.UpdatedAt
	}
	return l.svcCtx.AppUserProfileModel.Update(l.ctx, &appuser, nil)
}
