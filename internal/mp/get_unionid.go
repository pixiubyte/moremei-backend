package mp

import (
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/jsonx"
	"moremei/ai-saas/internal/mp/types"
	"moremei/ai-saas/pkg/httpc"
	"net/http"
)

const UnionIDURl = "https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN"

func (c *Client) GetUnionID(openid string) (string, error) {
	token, err := c.GetAccessToken()
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf(UnionIDURl, token, openid)
	body, err := httpc.Request(http.MethodGet, url, nil, nil)
	if err != nil {
		return "", err
	}

	var resp types.UnionIDResponse
	err = jsonx.UnmarshalFromString(body, &resp)
	if err != nil {
		return "", err
	}

	if resp.ErrCode != 0 {
		return "", errors.New(resp.ErrMsg)
	}

	return resp.UnionId, nil
}
