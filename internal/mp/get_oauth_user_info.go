package mp

import (
	"errors"
	"fmt"
	"moremei/ai-saas/internal/mp/types"
	"net/http"

	"moremei/ai-saas/pkg/httpc"

	"github.com/zeromicro/go-zero/core/jsonx"
)

const (
	OauthUrl = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	UserInfo = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN"
)

func (c *Client) GetOauthUserInfo(code string) (*types.UserInfoResponse, error) {
	body, err := httpc.Request(http.MethodGet, fmt.Sprintf(OauthUrl, c.Appid, c.Secret, code), nil, nil)
	if err != nil {
		return nil, err
	}
	var oauthResp types.OauthOpenidResp
	err = jsonx.UnmarshalFromString(body, &oauthResp)
	if err != nil {
		return nil, err
	}
	if oauthResp.ErrCode != 0 {
		return nil, errors.New(oauthResp.ErrMsg)
	}
	userInfoBody, err := httpc.Request(http.MethodGet, fmt.Sprintf(UserInfo, oauthResp.AccessToken, oauthResp.Openid), nil, nil)
	if err != nil {
		return nil, err
	}
	var userInfoResp *types.UserInfoResponse
	err = jsonx.UnmarshalFromString(userInfoBody, &userInfoResp)
	if err != nil {
		return nil, err
	}
	if oauthResp.ErrCode != 0 {
		return nil, errors.New(oauthResp.ErrMsg)
	}
	return userInfoResp, nil
}
