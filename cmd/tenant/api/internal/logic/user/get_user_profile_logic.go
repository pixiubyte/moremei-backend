package user

import (
	"context"
	"encoding/json"
	"errors"

	"moremei/ai-saas/cmd/tenant/api/internal/consts"
	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	"moremei/ai-saas/cmd/tenant/api/internal/utils"
	"moremei/ai-saas/pkg/gormc"
	errorx "moremei/ai-saas/pkg/x/error"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type GetUserProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserProfileLogic {
	return &GetUserProfileLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetUserProfileLogic) GetUserProfile() (*types.UserProfileResponse, error) {
	authUser, err := utils.GetAuthUserCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	var (
		isBoundMobile bool
		mobile        string
		user          = new(struct {
			Uuid       string
			Nickname   string
			Avatar     string
			Gender     int64
			RegisterAt int64
			Address    string
			Birthday   string
			Height     uint64
		})
		exProfile map[string]string
	)

	err = mr.Finish(func() error {
		u, err := l.svcCtx.UserModel.FindOne(l.ctx, authUser.UserId)
		if err != nil {
			if !errors.Is(err, gormc.ErrNotFound) {
				return err
			}
			return errorx.UnauthorizedError
		}
		user.Uuid = u.Uuid
		user.Nickname = u.Nickname
		user.Avatar = u.Avatar
		user.Gender = u.Gender
		user.RegisterAt = u.RegisterAt.Unix()
		user.Address = u.Address
		user.Birthday = u.Birthday
		user.Height = u.Height
		return nil
	}, func() error {
		webUser, err := l.svcCtx.WebUserModel.FindOneByUserId(l.ctx, authUser.UserId)
		if err != nil {
			if !errors.Is(err, gormc.ErrNotFound) {
				return err
			}
			return nil
		}
		mobile = webUser.Mobile
		isBoundMobile = true
		return nil
	})
	if err != nil {
		return nil, err
	}

	aup, err := l.svcCtx.AppUserProfileModel.FindOneByAppUuidKey(l.ctx, consts.AppNameDefault, user.Uuid, consts.UserProfileKeyArchive)
	if err != nil && !errors.Is(err, gormc.ErrNotFound) {
		return nil, err
	}
	if aup != nil {
		if err := json.Unmarshal([]byte(aup.Value), &exProfile); err != nil {
			return nil, err
		}
	}
	if exProfile == nil {
		exProfile = map[string]string{}
	}

	return &types.UserProfileResponse{
		Id:            user.Uuid,
		Nickname:      user.Nickname,
		Mobile:        mobile,
		Avatar:        user.Avatar,
		Gender:        int(user.Gender),
		RegistedTime:  user.RegisterAt,
		IsBoundMobile: isBoundMobile,
		IsBoundWechat: false,
		BoundWechcat:  "",
		Address:       user.Address,
		Birthday:      user.Birthday,
		Province:      exProfile["province"],
		City:          exProfile["city"],
		Subcompany:    exProfile["subcompany"],
		Height:        int64(user.Height),
	}, nil
}
