package mp

import (
	"encoding/base64"
	"moremei/ai-saas/internal/mp/types"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Client struct {
	Appid           string
	Secret          string
	Token           string
	AesKey          []byte
	AuthorizerAppid string
	Redis           *redis.Redis
}

type MpHandler interface {
	GetOauthUserInfo(code string) (*types.UserInfoResponse, error)
	GetAccessToken() (string, error)
	SendTemplateMsg(openid, templateId, appid, pagePath string, data interface{}) error
	HandleAccess(w http.ResponseWriter, r *http.Request)
	CreateTemporaryQRCodeTicket(sceneStr string, expireSeconds ...int) (ticket *QRCodeTicket, err error)
	GetUnionID(openid string) (string, error)
}

func NewMp(config Config, redis *redis.Redis) MpHandler {
	EncodingAESKey, err := base64.StdEncoding.DecodeString(config.AesKey + "=")
	if err != nil {
		logx.Errorf("aesKey 配置错误：%v", err)
	}
	return &Client{
		Appid:  config.Appid,
		Secret: config.Secret,
		Redis:  redis,
		Token:  config.Token,
		AesKey: EncodingAESKey,
	}
}
