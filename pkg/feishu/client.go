package feishu

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/jsonx"
	"moremei/ai-saas/pkg/httpc"
	"net/http"
)

const (
	CreateReqUrl  = "https://open.feishu.cn/open-apis/sheets/v3/spreadsheets"
	QuerySheetUrl = "https://open.feishu.cn/open-apis/sheets/v3/spreadsheets/%s/sheets/query"
	WriteSheetUrl = "https://open.feishu.cn/open-apis/sheets/v2/spreadsheets/%s/values_prepend"
)

type Client struct {
	appId  string
	secret string
}

func NewClient(appId, secret string) *Client {
	return &Client{
		appId:  appId,
		secret: secret,
	}
}

func (c *Client) getTenantAccessToken() (string, error) {
	req := fmt.Sprintf("{\"app_id\": \"%s\",\"app_secret\": \"%s\"}", c.appId, c.secret)
	body, err := httpc.Request(http.MethodPost, "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal",
		map[string]string{"Content-Type": "application/json; charset=utf-8"}, []byte(req))

	if err != nil {
		return "", err
	}
	var resp TokenResp
	err = jsonx.UnmarshalFromString(body, &resp)
	if err != nil {
		return "", err
	}
	if resp.Code != 0 {
		return "", errors.New(resp.Msg)
	}
	return resp.TenantAccessToken, nil
}

func (c *Client) CreateSheet(title string) (*Spreadsheet, error) {
	token, err := c.getTenantAccessToken()
	if err != nil {
		return nil, err
	}
	header := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
		"Content-Type":  "application/json",
	}
	req := SheetCreateReq{
		Title: title,
	}
	b, _ := json.Marshal(req)
	body, err := httpc.Request(http.MethodPost, CreateReqUrl, header, b)
	if err != nil {
		return nil, err
	}
	var resp SheetCreateResp
	err = jsonx.UnmarshalFromString(body, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, errors.New(resp.Msg)
	}
	return &resp.Data.Spreadsheet, nil
}

func (c *Client) QuerySheet(token string) (string, error) {
	url := fmt.Sprintf(QuerySheetUrl, token)

	token, err := c.getTenantAccessToken()
	if err != nil {
		return "", err
	}
	header := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}
	body, err := httpc.Request(http.MethodGet, url, header, nil)
	if err != nil {
		return "", err
	}
	var resp QuerySheetResp
	err = jsonx.UnmarshalFromString(body, &resp)
	if err != nil {
		return "", err
	}
	if resp.Code != 0 {
		return "", errors.New(resp.Msg)
	}
	return resp.Data.Sheets[0].SheetId, nil
}

func (c *Client) WriteSheet(token, sheetId string, data [][]interface{}) error {
	url := fmt.Sprintf(WriteSheetUrl, token)
	token, err := c.getTenantAccessToken()
	if err != nil {
		return err
	}
	header := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
		"Content-Type":  "application/json; charset=utf-8",
	}
	var req = WriteSheetReq{
		ValueRange{
			Range:  fmt.Sprintf("%s!%s:%s", sheetId, "A1", "C"),
			Values: data,
		},
	}
	b, _ := jsonx.Marshal(req)
	body, err := httpc.Request(http.MethodPost, url, header, b)
	if err != nil {
		return err
	}
	var resp WriteSheetResp
	err = jsonx.UnmarshalFromString(body, &resp)
	if err != nil {
		return err
	}
	if resp.Code != 0 {
		return errors.New(resp.Msg)
	}
	return nil
}
