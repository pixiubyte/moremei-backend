package miniprogram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"moremei/ai-saas/internal/miniprogram/types"
	"moremei/ai-saas/pkg/httpc"

	"github.com/zeromicro/go-zero/core/jsonx"
)

func (d *defaultMiniProgram) GetUrlLink(req *types.GetUrlLinkRequest) (string, error) {
	token, err := d.GetAccessToken()
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false) // 关键步骤：禁用HTML符号转义
	err = encoder.Encode(req)

	if err != nil {
		return "", err
	}

	url := fmt.Sprintf(GetUrlLinkUrl, token)
	resp, err := httpc.Request(http.MethodPost, url, nil, buffer.Bytes())
	if err != nil {
		return "", err
	}

	var respBody types.GetUrlLinkResponse
	if err := jsonx.UnmarshalFromString(resp, &respBody); err != nil {
		return "", err
	}
	if respBody.Errcode != 0 {
		return "", fmt.Errorf("get url link failed, errcode: %d, errmsg: %s", respBody.Errcode, respBody.Errmsg)
	}

	return respBody.UrlLink, nil
}
