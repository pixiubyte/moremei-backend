package miniprogram

import (
	"moremei/ai-saas/internal/miniprogram/types"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	JsCode2SessionUrl     = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	UserPhoneNumberUrl    = "https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=%s"
	AccessTokenUrl        = "https://api.weixin.qq.com/cgi-bin/stable_token"
	SendMessageUrl        = "https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token=%s"
	GetUnlimitedQRCodeUrl = "https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=%s"
	GetUrlLinkUrl         = "https://api.weixin.qq.com/wxa/generate_urllink?access_token=%s"
)

var _ MiniProgram = (*defaultMiniProgram)(nil)

type MiniProgram interface {
	GetMobile(code string) (string, error)
	GetOpenid(code string) (*types.GetOpenidResponse, error)
	SendMessage(toUser, templateId, url string, data interface{}) error
	GetUnlimitedQRCode(req *types.GetUnlimitedQRCodeRequest) ([]byte, error)
	GetUrlLink(req *types.GetUrlLinkRequest) (string, error)
}

type defaultMiniProgram struct {
	AppId     string
	AppSecret string
	Redis     *redis.Redis
	Mod       string
}

func NewMiniProgram(config Config, redis *redis.Redis, mod string) MiniProgram {
	return &defaultMiniProgram{
		AppId:     config.AppId,
		AppSecret: config.AppSecret,
		Redis:     redis,
		Mod:       mod,
	}
}
