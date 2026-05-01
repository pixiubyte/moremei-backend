package mp

import (
	"encoding/json"
	"fmt"
	"moremei/ai-saas/internal/mp/types"
	"net/http"

	"moremei/ai-saas/pkg/consts"
	"moremei/ai-saas/pkg/httpc"

	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/logx"
)

const AccessTokenUrl = "https://api.weixin.qq.com/cgi-bin/stable_token"

func (c *Client) GetAccessToken() (string, error) {
	key := fmt.Sprintf(consts.AccessToken, "mp", "mei", c.Appid)
	val, err := c.Redis.Get(key)
	if err != nil {
		return "", err
	}

	if val != "" {
		return val, nil
	} else {
		request := types.AccessTokenRequest{
			GrantType: "client_credential",
			Appid:     c.Appid,
			Secret:    c.Secret,
		}
		bytes, err := json.Marshal(request)
		if err != nil {
			return "", err
		}
		body, err := httpc.Request(http.MethodPost, AccessTokenUrl, nil, bytes)
		if err != nil {
			return "", err
		}

		var tokenResp types.AccessTokenResp
		err = jsonx.UnmarshalFromString(body, &tokenResp)
		if err != nil {
			return "", err
		}

		if tokenResp.ErrCode != 0 {
			return "", fmt.Errorf("%s", tokenResp.ErrMsg)
		}

		_, err = c.Redis.SetnxEx(key, tokenResp.AccessToken, tokenResp.ExpiresIn-150)
		if err != nil {
			logx.Errorf("redis set err:%v", err)
		}

		return tokenResp.AccessToken, nil
	}
}
