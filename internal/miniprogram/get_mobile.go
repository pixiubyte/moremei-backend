package miniprogram

import (
	"errors"
	"fmt"
	"net/http"

	"moremei/ai-saas/internal/miniprogram/types"
	"moremei/ai-saas/pkg/httpc"

	"github.com/zeromicro/go-zero/core/jsonx"
)

func (c *defaultMiniProgram) GetMobile(code string) (string, error) {
	request := fmt.Sprintf(`{"code":"%s"}`, code)
	token, err := c.GetAccessToken()
	if err != nil {
		return "", err
	}
	body, err := httpc.Request(http.MethodPost, fmt.Sprintf(UserPhoneNumberUrl, token), nil, []byte(request))
	if err != nil {
		return "", err
	}
	var resp types.GetMobileResponse
	err = jsonx.UnmarshalFromString(body, &resp)
	if err != nil {
		return "", err
	}
	if resp.Errcode != 0 {
		return "", errors.New(resp.Errmsg)
	}
	return resp.PhoneInfo.PhoneNumber, nil
}
