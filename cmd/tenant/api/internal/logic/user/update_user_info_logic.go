package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"moremei/ai-saas/cmd/tenant/api/internal/consts"
	"moremei/ai-saas/cmd/tenant/api/internal/errs"
	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	"moremei/ai-saas/cmd/tenant/api/internal/utils"
	"moremei/ai-saas/internal/model/tenant"
	"moremei/ai-saas/pkg/gormc"
	pkgUtils "moremei/ai-saas/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type UpdateUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(req *types.UpdateUserInfo) error {
	authUser, err := utils.GetAuthUserCtx(l.ctx)
	if err != nil {
		return err
	}
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, authUser.UserId)
	if err != nil {
		return err
	}
	if req.Nickname != "" {
		if len(req.Nickname) > 30 {
			return errs.UserNameLongError
		}
		user.Nickname = req.Nickname
	}
	if req.Sex != 0 {
		user.Gender = req.Sex
	}
	if req.Birthday != "" {
		if _, err := time.Parse(time.DateOnly, req.Birthday); err != nil {
			return errs.BirthDayError
		}
		user.Birthday = req.Birthday
	}
	if req.Address != "" {
		if len(req.Address) > 255 {
			return errs.AddressLongError
		}
		user.Address = req.Address
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if err := l.svcCtx.UserModel.Update(l.ctx, user, nil); err != nil {
		return err
	}

	userAidLock := fmt.Sprintf("%s:%s:%v:%v", l.svcCtx.Config.Name, consts.AppNameDefault, consts.UserProfileKeyArchive, authUser.UserId)
	lock := redis.NewRedisLock(l.svcCtx.RedisClient, userAidLock)
	if err := pkgUtils.AcquireAndWaitRedisLock(l.ctx, lock, 5); err != nil {
		return err
	}
	defer lock.ReleaseCtx(l.ctx)

	usf, err := l.svcCtx.AppUserProfileModel.FindOneByAppUuidKey(l.ctx, consts.AppNameDefault, user.Uuid, consts.UserProfileKeyArchive)
	if err != nil {
		if !errors.Is(err, gormc.ErrNotFound) {
			return err
		}
		usf = &tenant.AppUserProfile{App: consts.AppNameDefault, Uuid: user.Uuid, Key: consts.UserProfileKeyArchive, Value: "{}"}
	}

	exProfile := map[string]string{}
	if usf.Value != "" {
		if err := json.Unmarshal([]byte(usf.Value), &exProfile); err != nil {
			return err
		}
	}
	if req.Province != "" {
		exProfile["province"] = req.Province
	}
	if req.City != "" {
		exProfile["city"] = req.City
	}
	if req.Subcompany != "" {
		exProfile["subcompany"] = req.Subcompany
	}
	exValue, err := json.Marshal(exProfile)
	if err != nil {
		return err
	}
	usf.Value = string(exValue)
	return l.svcCtx.AppUserProfileModel.Update(l.ctx, usf, nil)
}
