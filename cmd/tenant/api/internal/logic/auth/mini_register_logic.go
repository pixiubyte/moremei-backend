package auth

import (
	"context"
	"time"
	"github.com/zeromicro/go-zero/core/jsonx"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"moremei/ai-saas/cmd/tenant/api/internal/consts"
	"moremei/ai-saas/cmd/tenant/api/internal/errs"
	"moremei/ai-saas/cmd/tenant/api/internal/logic/sms"
	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	"moremei/ai-saas/internal/model/tenant"
	pkgconsts "moremei/ai-saas/pkg/consts"
	"moremei/ai-saas/pkg/gormc"
	"moremei/ai-saas/pkg/jwt"
	errorx "moremei/ai-saas/pkg/x/error"
)

type MiniRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMiniRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MiniRegisterLogic {
	return &MiniRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MiniRegisterLogic) MiniRegister(req *types.MiniRegisterRequest) (resp *types.LoginResponse, err error) {
	mobile := req.Mobile
	miniClient, err := l.svcCtx.SdkClients.GetMiniClient(l.ctx)
	if err != nil {
		return nil, err
	}
	// 获取用户openid
	respOpenid, err := miniClient.GetOpenid(req.JsCode)
	if err != nil {
		l.Logger.Errorf("获取openid失败err:%v", err)
		return nil, err
	}
	// 手机号码验证
	err = sms.NewSmsLogic(l.ctx, l.svcCtx).Verify(mobile, req.Code, tenant.SmsBizTypeBind)
	if err != nil {
		return nil, err
	}

	var (
		user     *tenant.User     // 用户
		webUser  *tenant.WebUser  // 用户关联手机号
		miniUser *tenant.MiniUser // 用户关联小程序用户
	)

	miniUser, err = l.svcCtx.MiniUserModel.FindOneByOpenid(l.ctx, respOpenid.Openid)
	//存在小程序用户
	if err == nil {
		_, err = l.svcCtx.WebUserModel.FindOneByUserId(l.ctx, miniUser.UserId)
		if err == nil {
			return nil, errorx.NewError(2004, "请勿重复绑定")
		} else if !errors.Is(err, gormc.ErrNotFound) {
			return nil, errors.Errorf("user:%d, find web user err:%v", miniUser.UserId, err)
		}
		user, err = l.svcCtx.UserModel.FindOne(l.ctx, miniUser.UserId)
		if err != nil {
			return nil, errors.Errorf("user:%d, find user err:%v", miniUser.UserId, err)
		}
		// 禁用状态用户
		if pkgconsts.IsFlagFalse(user.Status) {
			return nil, errorx.ForbiddenError
		}
		webUser = &tenant.WebUser{
			UserId: miniUser.UserId,
			Mobile: mobile,
		}
		err = l.svcCtx.WebUserModel.Insert(l.ctx, webUser, nil)
		if err != nil {
			return nil, errors.Errorf("user:%d, insert web user err:%v", webUser.UserId, err)
		}
		return l.login(user)
	} else if !errors.Is(err, gormc.ErrNotFound) {
		return nil, errors.Errorf("openid:%s, find mini user err:%v", respOpenid.Openid, err)
	}

	miniUser = &tenant.MiniUser{
		Openid:  respOpenid.Openid,
		Unionid: respOpenid.Unionid,
	}
	webUser, err = l.svcCtx.WebUserModel.FindOneByMobile(l.ctx, mobile)

	nickname, _ := l.hasEihConfig(mobile)
	// 已存在web用户
	if err == nil {
		_, err = l.svcCtx.MiniUserModel.FindOneByUserId(l.ctx, webUser.UserId)
		if err == nil {
			return nil, errs.MiniBindError
		} else if !errors.Is(err, gormc.ErrNotFound) {
			return nil, errors.Errorf("user:%d, find exist mini user err:%v", webUser.UserId, err)
		}

		user, err = l.svcCtx.UserModel.FindOne(l.ctx, webUser.UserId)
		if err != nil {
			return nil, errors.Errorf("user:%d, find user err:%v", webUser.UserId, err)
		}
		// 禁用状态用户
		if pkgconsts.IsFlagFalse(user.Status) {
			return nil, errorx.ForbiddenError
		}

		// 新增小程序用户
		miniUser.UserId = webUser.UserId
		err = l.svcCtx.MiniUserModel.Insert(l.ctx, miniUser, nil)
		if err != nil {
			return nil, errors.Errorf("user:%d, insert mini user err:%v", miniUser.UserId, err)
		}
		// 登录
		return l.login(user)
	} else if !errors.Is(err, gormc.ErrNotFound) {
		return nil, errors.Errorf("mobile:%s, find web user err:%v", mobile, err)
	}

	// 新注册用户
	user = &tenant.User{
		Uuid:       l.svcCtx.Snowflake.Generate(),
		Nickname:   nickname,
		Avatar:     "",
		RegisterAt: time.Now(),
		Status:     pkgconsts.FlagTrue,
	}
	webUser = &tenant.WebUser{
		Mobile: mobile,
	}
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		err = l.svcCtx.UserModel.Insert(l.ctx, user, tx)
		if err != nil {
			return err
		}
		webUser.UserId = user.Id
		miniUser.UserId = user.Id
		err = l.svcCtx.WebUserModel.Insert(l.ctx, webUser, tx)
		if err != nil {
			return err
		}

		return l.svcCtx.MiniUserModel.Insert(l.ctx, miniUser, tx)
	})
	if err != nil {
		logx.Errorf("mobile:%s, insert user err:%v", mobile, err)
		return nil, err
	}

	// 登录
	return l.login(user)
}

func (l *MiniRegisterLogic) login(user *tenant.User) (resp *types.LoginResponse, err error) {
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
func (l *MiniRegisterLogic) hasEihConfig(mobile string) (string, bool) {
	// 获取手机号后4位
	suffix := mobile
	if len(mobile) > 4 {
		suffix = mobile[len(mobile)-4:]
	}
	nickname := "酶好用户" + suffix

	// 获取系统配置
	systemConfig, err := l.svcCtx.SystemConfigModel.FindByKey(l.ctx, tenant.SystemConfigKeyEihDefaultConfig)
	if err != nil {
		if errors.Is(err, gormc.ErrNotFound) {
			return nickname, false
		}
		logx.Errorf("hasEihConfig mobile:%s, find system config err:%v", mobile, err)
		return nickname, false
	}

	// 检查配置是否为空
	if systemConfig == nil || systemConfig.Value == "" {
		logx.Infof("EIH配置为空，使用默认昵称")
		return nickname, false
	}

	// 解析配置
	eihDefaultConfig := new(tenant.EihDefaultConfig)
	err = jsonx.UnmarshalFromString(systemConfig.Value, &eihDefaultConfig)
	if err != nil {
		logx.Errorf("hasEihConfig mobile:%s, unmarshal system config err:%v, raw value:%s",
			mobile, err, systemConfig.Value)
		return nickname, false
	}

	// 检查配置对象是否为空
	if eihDefaultConfig == nil {
		logx.Errorf("EIH配置解析后为空")
		return nickname, false
	}

	// 使用配置的昵称前缀
	if eihDefaultConfig.DefaultNickname != "" {
		nickname = eihDefaultConfig.DefaultNickname + suffix
	}

	return nickname, eihDefaultConfig.HasHabit == 1
}
