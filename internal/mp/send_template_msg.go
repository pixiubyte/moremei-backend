package mp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"moremei/ai-saas/internal/mp/types"
	"moremei/ai-saas/pkg/httpc"

	"github.com/zeromicro/go-zero/core/jsonx"
)

const TemplateUrl = "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s"

func (c *Client) SendTemplateMsg(openid, templateId, appid, pagePath string, data interface{}) error {
	accessToken, err := c.GetAccessToken()
	if err != nil {
		return err
	}
	templateReq := types.Template{
		Touser:     openid,
		TemplateId: templateId,
		Miniprogram: types.Miniprogram{
			Appid:    appid,
			Pagepath: pagePath,
		},
		Data: data,
	}
	reqBody, err := json.Marshal(templateReq)
	if err != nil {
		return err
	}
	body, err := httpc.Request(http.MethodPost, fmt.Sprintf(TemplateUrl, accessToken), nil, reqBody)
	if err != nil {
		return err
	}
	var common types.Common
	err = jsonx.UnmarshalFromString(body, &common)
	if err != nil {
		return err
	}
	if common.ErrCode != 0 {
		return errors.New(common.ErrMsg)
	}
	return nil
}
