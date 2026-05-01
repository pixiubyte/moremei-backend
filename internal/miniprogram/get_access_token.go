package miniprogram

import (
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/jsonx"
	"moremei/ai-saas/internal/miniprogram/types"
	"moremei/ai-saas/pkg/httpc"
	"net/http"
)

func (c *defaultMiniProgram) GetAccessToken() (string, error) {
	key := fmt.Sprintf("miniprogram:%s:%s", c.AppId, c.AppSecret)
	val, err := c.Redis.Get(key)
	if err != nil {
		return "", err
	}
	if val != "" {
		return val, nil
	}

	request := types.StableToken{
		GrantType:    "client_credential",
		Appid:        c.AppId,
		Secret:       c.AppSecret,
		ForceRefresh: false,
	}

	bytes, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	body, err := httpc.Request(http.MethodPost, AccessTokenUrl, nil, bytes)
	if err != nil {
		return "", err
	}

	var resp types.AccessTokenResponse
	err = jsonx.UnmarshalFromString(body, &resp)
	if err != nil {
		return "", err
	}
	if resp.Errcode != 0 {
		return "", fmt.Errorf("%d：%s", resp.Errcode, resp.Errmsg)
	}

	err = c.Redis.Setex(key, resp.AccessToken, resp.ExpiresIn-200)
	if err != nil {
		return "", err
	}
	return resp.AccessToken, nil
}
