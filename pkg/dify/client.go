package dify

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	host         string
	apiSecretKey string

	httpClient  *http.Client
	httpRequest *http.Request

	thinking bool
}

type Option func(c *Client)

type ErrResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func NewClientWithConfig(c *ClientConfig) *Client {
	var httpClient = &http.Client{}

	if c.Timeout == 0 {
		httpClient.Timeout = c.Timeout
	}
	if c.Transport != nil {
		httpClient.Transport = c.Transport
	}

	return &Client{
		host:         c.Host,
		apiSecretKey: c.ApiSecretKey,
		httpClient:   httpClient,
	}
}

func NewClient(host, apiSecretKey string, options ...Option) *Client {
	options = append(options, WithHost(host), WithApiSecretKey(apiSecretKey))
	c := &Client{
		httpClient: &http.Client{Timeout: 0},
	}
	for _, option := range options {
		option(c)
	}
	return c
}

func WithHost(host string) Option {
	return func(c *Client) {
		c.host = host
	}
}

func WithApiSecretKey(apiSecretKey string) Option {
	return func(c *Client) {
		c.apiSecretKey = apiSecretKey
	}
}
func WithHttpClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}
func WithThinking(thinking bool) Option {
	return func(c *Client) {
		c.thinking = thinking
	}
}

func (c *Client) NewHttpRequest(ctx context.Context, method, requestUrl string, request ...interface{}) (r *http.Request, err error) {
	if method == http.MethodGet {
		if len(request) > 0 {
			if urlValues, ok := request[0].(url.Values); ok {
				var requestUrlParse *url.URL
				if requestUrlParse, err = url.Parse(requestUrl); err != nil {
					return
				}
				requestUrlParse.RawQuery = urlValues.Encode()
				requestUrl = requestUrlParse.String()
			}
		}
		r, err = http.NewRequestWithContext(ctx, method, requestUrl, http.NoBody)
	} else if method == http.MethodPost {
		var b io.Reader
		if len(request) > 0 {
			var reqBytes []byte
			if reqBytes, err = json.Marshal(request[0]); err != nil {
				return
			}
			b = bytes.NewBuffer(reqBytes)
		} else {
			b = http.NoBody
		}
		r, err = http.NewRequestWithContext(ctx, method, requestUrl, b)
	} else {
		err = errors.New("NewHttpRequest.method must be http.MethodGet or http.MethodPost")
	}
	return
}

func (c *Client) SetHttpRequest(r *http.Request) *Client {
	c.httpRequest = r
	return c
}

func (c *Client) SetHttpRequestHeader(key string, value string) *Client {
	c.httpRequest.Header.Set(key, value)
	return c
}

func (c *Client) SendRequest(res interface{}) (err error) {
	if c.httpRequest == nil {
		panic("http_request illegal")
	}
	var resp *http.Response
	if resp, err = c.httpClient.Do(c.httpRequest); err != nil {
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var errRes ErrResponse
		if err = json.NewDecoder(resp.Body).Decode(&errRes); err != nil {
			return err
		}
		err = errors.New(errRes.Message)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(res)
	return
}

func (c *Client) SendRequestStream() (resp *http.Response, err error) {
	if c.httpRequest == nil {
		panic("http_request illegal")
	}
	resp, err = c.httpClient.Do(c.httpRequest)
	return
}

func (c *Client) GetHost() string {
	var host = strings.TrimSuffix(c.host, "/")
	return host
}

func (c *Client) GetApiSecretKey() string {
	return c.apiSecretKey
}

func (c *Client) buildRequestApi(requestUrl string) string {
	return c.GetHost() + requestUrl
}

func (c *Client) createBaseRequest(ctx context.Context, method string, url string, req ...interface{}) (r *http.Request, err error) {
	if r, err = c.NewHttpRequest(ctx, method, url, req...); err == nil {
		c.SetHttpRequest(r).
			SetHttpRequestHeader("Authorization", "Bearer "+c.GetApiSecretKey()).
			SetHttpRequestHeader("Cache-Control", "no-cache").
			SetHttpRequestHeader("Content-Type", "application/json; charset=utf-8")
	}
	return
}

func (c *Client) createGetRequest(ctx context.Context, url string, req ...interface{}) (*http.Request, error) {
	return c.createBaseRequest(ctx, http.MethodGet, url, req...)
}

func (c *Client) createPostRequest(ctx context.Context, url string, req ...interface{}) (*http.Request, error) {
	return c.createBaseRequest(ctx, http.MethodPost, url, req...)
}

func (c *Client) trimThinkingContent(content string) string {
	// 以下方法测试后无法去除 thinking 内容
	//re := regexp.MustCompile(`<details.*>.*</details>(.*)`)
	//tmpContent, err := strconv.Unquote(strings.Replace(strconv.Quote(content), `\\u`, `\u`, -1))
	//if err != nil {
	//	return content
	//}
	//match := re.FindStringSubmatch(tmpContent)
	//if len(match) == 2 {
	//	content = strings.Trim(match[1], "\n")
	//}
	substrs := []string{"</details>", "</think>"}
	for _, substr := range substrs {
		li := strings.LastIndex(content, substr)
		if li > 0 && len(content) > li+len(substr) {
			content = strings.Trim(content[li+len(substr):], "\n")
		}
	}

	return content
}
