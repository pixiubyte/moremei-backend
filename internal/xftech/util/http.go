package util

import (
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

type CommonResponse struct {
	Code    int         `json:"Code"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
	Time    int64       `json:"t"`
}

func GetHttpClient() *resty.Client {
	return resty.New().SetHeaders(map[string]string{
		"Content-Type": "application/json; charset=utf-8",
		"Accept":       "application/json",
	}).SetTimeout(0)
}

func HttpGet(url string) (resp *CommonResponse, err error) {
	c := GetHttpClient()

	//var response CommonResponse
	response, err := c.R().SetResult(&CommonResponse{}).Get(url)
	if err != nil {
		return nil, err
	}

	if !response.IsSuccess() {
		return nil, errors.Errorf("请求出错，code：%d，body：%s", response.StatusCode(), response.Body())
	}

	resp, ok := response.Result().(*CommonResponse)
	if !ok {
		return nil, errors.New("返回数据无法解析")
	}

	return resp, nil
}

func HttpPostWithFormData(url string, formData map[string]string) (resp *CommonResponse, err error) {
	c := GetHttpClient().
		SetHeader("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	response, err := c.R().SetFormData(formData).SetResult(&CommonResponse{}).Post(url)
	if err != nil {
		return nil, err
	}

	resp, ok := response.Result().(*CommonResponse)
	if !ok {
		return nil, errors.New("返回数据无法解析")
	}

	return resp, nil
}
