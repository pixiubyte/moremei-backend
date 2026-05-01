package miniprogram

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"moremei/ai-saas/internal/miniprogram/types"
	"moremei/ai-saas/pkg/httpc"

	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/service"
)

func (c *defaultMiniProgram) SendMessage(toUser, templateId, url string, data interface{}) error {
	request := types.SendMessageRequest{
		Touser:     toUser,
		TemplateId: templateId,
		Page:       url,
		Data:       data,
		Lang:       "zh_CN",
	}
	if c.Mod == service.DevMode {
		request.MiniprogramState = "trial"
	} else {
		request.MiniprogramState = "formal"
	}
	token, err := c.GetAccessToken()
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(request)
	if err != nil {
		return err
	}
	body, err := httpc.Request(http.MethodPost, fmt.Sprintf(SendMessageUrl, token), nil, bytes)
	if err != nil {
		return err
	}

	var resp types.CommonResponse
	err = jsonx.UnmarshalFromString(body, &resp)
	if err != nil {
		return err
	}

	if resp.Errcode != 0 && resp.Errcode != 43101 {
		return errors.New(resp.Errmsg)
	}
	return nil
}
