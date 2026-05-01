package ainlp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/samber/lo"
	"io"
	"net/http"
	"strings"
)

type (
	Client struct {
		config     *Config
		httpClient *http.Client
	}

	Option func(*Client)

	CommonResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    any    `json:"data"`
		Time    int64  `json:"time"`
	}

	ExtractKeywordsRequest struct {
		Content string `json:"content"`
	}
)

const SuccessCode = 0

func NewClient(cfg *Config, opts ...Option) *Client {
	c := &Client{config: cfg, httpClient: &http.Client{}}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithHttpClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func (c *Client) ExtractKeywords(ctx context.Context, text string) ([]string, error) {
	req, err := c.newRequest(ctx, http.MethodPost, "v1/keyphrase", &ExtractKeywordsRequest{Content: text})
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	commResp := CommonResponse{}
	err = json.Unmarshal(response, &commResp)
	if err != nil {
		return nil, err
	}

	if commResp.Code != SuccessCode {
		return nil, errors.New(commResp.Message)
	}

	if commResp.Data == nil {
		return nil, errors.New("data is nil")
	}

	data, ok := commResp.Data.([]interface{})
	if !ok {
		return nil, errors.New("data is not []interface{}")
	}

	return lo.Map[interface{}, string](data, func(item interface{}, index int) string {
		return item.(string)
	}), nil
}

func (c *Client) newRequest(
	ctx context.Context,
	method string,
	urlSuffix string,
	body interface{}) (*http.Request, error) {
	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, formatURL(c.config.BaseUrl, urlSuffix), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	return req, nil
}

func formatURL(url, suffix string) string {
	if !strings.HasSuffix(url, suffix) {
		if strings.LastIndex(url, "/") != len(url)-1 {
			url += "/"
		}
		url += suffix
	}
	return url
}
