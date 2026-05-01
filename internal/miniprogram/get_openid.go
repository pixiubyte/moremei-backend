package miniprogram

import (
	"errors"
	"fmt"
	"net/http"

	"moremei/ai-saas/internal/miniprogram/types"
	"moremei/ai-saas/pkg/httpc"

	"github.com/zeromicro/go-zero/core/jsonx"
)

func (c *defaultMiniProgram) GetOpenid(code string) (*types.GetOpenidResponse, error) {
	body, err := httpc.Request(http.MethodGet, fmt.Sprintf(JsCode2SessionUrl, c.AppId, c.AppSecret, code), nil, nil)
	if err != nil {
		return nil, err
	}
	var resp *types.GetOpenidResponse
	err = jsonx.UnmarshalFromString(body, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Errcode != 0 {
		return nil, errors.New(resp.Errmsg)
	}
	return resp, nil
}
