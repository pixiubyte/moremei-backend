package mp

import (
	"fmt"
	"net/http"
	"net/url"

	"moremei/ai-saas/internal/mp/types"
	"moremei/ai-saas/pkg/httpc"

	"github.com/zeromicro/go-zero/core/jsonx"
)

// 二维码
const (
	AccountCreateQRCodeURL     = "https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=%s"
	AccountGetQRCodeImgAddrURL = "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=%s"
)

// QRCodeTicket 二维码ticket
type QRCodeTicket struct {
	types.Common
	Ticket        string `json:"ticket"`         // 获取的二维码ticket，凭借此ticket可以在有效时间内换取二维码。
	ExpireSeconds int    `json:"expire_seconds"` // 该二维码有效时间，以秒为单位。 最大不超过2592000（即30天）。
	URL           string `json:"url"`            // 二维码图片解析后的地址，开发者可根据该地址自行生成需要的二维码图片
}

// CreateTemporaryQRCodeTicket 创建临时二维码
func (c *Client) CreateTemporaryQRCodeTicket(sceneStr string, expireSeconds ...int) (ticket *QRCodeTicket, err error) {
	if len(expireSeconds) == 0 {
		expireSeconds = []int{30}
	} else if expireSeconds[0] >= 2592000 || expireSeconds[0] < 0 {
		expireSeconds[0] = 2592000
	}
	body := fmt.Sprintf(`{"expire_seconds": %d, "action_name": "QR_STR_SCENE", "action_info": {"scene": {"scene_str": "%s"}}}`, expireSeconds[0], sceneStr)
	return c.createQRCodeTicket(body)
}

func (c *Client) createQRCodeTicket(js string) (ticket *QRCodeTicket, err error) {
	ticket = &QRCodeTicket{}
	token, err := c.GetAccessToken()
	if err != nil {
		return nil, err
	}
	body, err := httpc.Request(http.MethodPost, fmt.Sprintf(AccountCreateQRCodeURL, token), nil, []byte(js))
	if err != nil {
		return nil, err
	}
	err = jsonx.UnmarshalFromString(body, ticket)
	if err != nil {
		return nil, err
	}
	return ticket, err
}

// GetQRCodeImg 通过ticket换取二维码
func GetQRCodeImg(ticket string) string {
	return fmt.Sprintf(AccountGetQRCodeImgAddrURL, url.QueryEscape(ticket))
}
