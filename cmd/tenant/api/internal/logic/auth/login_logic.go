package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"moremei/ai-saas/cmd/tenant/api/internal/consts"
	"moremei/ai-saas/cmd/tenant/api/internal/logic/sms"
	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	"moremei/ai-saas/internal/model/tenant"
	commonConst "moremei/ai-saas/pkg/consts"
	"moremei/ai-saas/pkg/gormc"
	"moremei/ai-saas/pkg/jwt"
	"moremei/ai-saas/pkg/uuid"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"gorm.io/gorm"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

var verifyMobileScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1], KEYS[2])
else
    return 0
end`

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	if l.svcCtx.Config.Mode != service.DevMode {
		//校验短信验证码
		err = sms.NewSmsLogic(l.ctx, l.svcCtx).Verify(req.Mobile, req.Code, tenant.SmsBizTypeLogin)
		if err != nil {
			return nil, err
		}
	}

	var (
		user    *tenant.User
		webUser *tenant.WebUser
	)
	webUser, err = l.svcCtx.WebUserModel.FindOneByMobile(l.ctx, req.Mobile)
	if err != nil {
		if !errors.Is(err, gormc.ErrNotFound) {
			return nil, err
		}
		//创建用户
		user, webUser, err = l.CreateUser(req.Mobile)
		if err != nil {
			return nil, err
		}
	} else {
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

	return &types.LoginResponse{
		AccessToken: token,
		ExpiresIn:   expiresTime,
		UserId:      user.Uuid,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
	}, nil
}

func (l *LoginLogic) CreateUser(mobile string) (*tenant.User, *tenant.WebUser, error) {
	// 生成昵称，随机6位字符串+手机号后4位
	nickname := fmt.Sprintf("%s%s", uuid.Random.GenerateN(6), mobile[len(mobile)-4:])
	user := tenant.User{
		Uuid:       l.svcCtx.Snowflake.Generate(),
		Nickname:   nickname,
		Avatar:     "",
		Gender:     tenant.UserGenderUnspecified,
		Status:     commonConst.FlagTrue,
		RegisterAt: time.Now(),
	}
	webUser := tenant.WebUser{
		Mobile:  mobile,
		Deleted: commonConst.FlagFalse,
	}
	aiRobot, err := l.svcCtx.AiRobotModel.FindDefaultOne(l.ctx)
	if err != nil {
		return nil, nil, err
	}
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		if err := l.svcCtx.UserModel.Insert(l.ctx, &user, tx); err != nil {
			return err
		}
		webUser.UserId = user.Id
		if err := l.svcCtx.WebUserModel.Insert(l.ctx, &webUser, tx); err != nil {
			return err
		}
		return l.svcCtx.ConversationModel.Insert(l.ctx, &tenant.Conversation{
			Uuid:      l.svcCtx.Snowflake.Generate(),
			UserId:    user.Id,
			AiRobotId: aiRobot.Id,
		}, tx)
	})
	if err != nil {
		return nil, nil, err
	}

	return &user, &webUser, nil
}
