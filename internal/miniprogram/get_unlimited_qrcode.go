package miniprogram

import (
	"encoding/json"
	"fmt"

	"moremei/ai-saas/internal/miniprogram/types"
	"moremei/ai-saas/pkg/httpc"

	"net/http"
)

func (c *defaultMiniProgram) GetUnlimitedQRCode(req *types.GetUnlimitedQRCodeRequest) ([]byte, error) {
	token, err := c.GetAccessToken()
	if err != nil {
		return nil, err
	}
	reqbody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	body, err := httpc.Request(http.MethodPost, fmt.Sprintf(GetUnlimitedQRCodeUrl, token), nil, reqbody)
	if err != nil {
		return nil, err
	}

	return []byte(body), nil
}
