package alisms

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type Client struct {
	config *Config
}

type Request struct {
	SignName     string
	TemplateId   string
	PhoneNumbers []string
	Params       map[string]interface{}
}

type Response struct {
	Result      bool
	PhoneNumber string
	RequestId   string
	Code        string
	Message     string
	BizId       string
}

type aliResponse struct {
	Message   string `json:"Message"`
	RequestId string `json:"RequestId"`
	Code      string `json:"Code"`
	BizId     string `json:"BizId"`
}

const (
	SendSuccessCode = "OK"
	endpoint        = "https://dysmsapi.aliyuncs.com"
	apiVersion      = "2017-05-25"
	signMethod      = "HMAC-SHA1"
	format          = "JSON"
	action          = "SendSms"
)

func MustNewClient(cfg *Config) *Client {
	return &Client{config: cfg}
}

func (c *Client) Send(req *Request) (map[string]Response, error) {
	err := c.validateRequest(req)
	if err != nil {
		return nil, err
	}

	// 将手机号列表转换为字符串
	phoneNumbers := strings.Join(req.PhoneNumbers, ",")

	// 将参数转换为JSON字符串
	templateParam := ""
	if len(req.Params) > 0 {
		paramBytes, err := json.Marshal(req.Params)
		if err != nil {
			return nil, errors.Wrap(err, "marshal template params failed")
		}
		templateParam = string(paramBytes)
	}

	// 准备请求参数
	params := map[string]string{
		"AccessKeyId":      c.config.AccessKeyId,
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"Format":           format,
		"SignatureMethod":  signMethod,
		"SignatureVersion": "1.0",
		"SignatureNonce":   fmt.Sprintf("%d", time.Now().UnixNano()),
		"Version":          apiVersion,
		"Action":           action,
		"RegionId":         c.config.Region,
		"PhoneNumbers":     phoneNumbers,
		"SignName":         req.SignName,
		"TemplateCode":     req.TemplateId,
		"TemplateParam":    templateParam,
	}

	// 计算签名
	signature := c.computeSignature(params)
	params["Signature"] = signature

	// 构建请求URL
	reqURL := endpoint + "/?" + c.buildQuery(params)

	// 发送请求
	resp, err := http.Get(reqURL)
	if err != nil {
		return lo.SliceToMap[string, string, Response](req.PhoneNumbers, func(item string) (string, Response) {
			return item, Response{
				PhoneNumber: item,
				Result:      false,
				Message:     err.Error(),
			}
		}), nil
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return lo.SliceToMap[string, string, Response](req.PhoneNumbers, func(item string) (string, Response) {
			return item, Response{
				PhoneNumber: item,
				Result:      false,
				Message:     err.Error(),
			}
		}), nil
	}

	// 解析响应
	var aliResp aliResponse
	err = json.Unmarshal(body, &aliResp)
	if err != nil {
		return lo.SliceToMap[string, string, Response](req.PhoneNumbers, func(item string) (string, Response) {
			return item, Response{
				PhoneNumber: item,
				Result:      false,
				Message:     err.Error(),
			}
		}), nil
	}

	// 为每个手机号生成结果
	result := Response{
		RequestId: aliResp.RequestId,
		Code:      aliResp.Code,
		Message:   aliResp.Message,
		BizId:     aliResp.BizId,
		Result:    aliResp.Code == SendSuccessCode,
	}

	return lo.SliceToMap[string, string, Response](req.PhoneNumbers, func(phoneNumber string) (string, Response) {
		result.PhoneNumber = phoneNumber
		return phoneNumber, result
	}), nil
}

func (c *Client) computeSignature(params map[string]string) string {
	// 1. 按参数名称字典序排序
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 2. 构建规范化请求字符串
	var canonicalizedQueryString strings.Builder
	for i, k := range keys {
		if i > 0 {
			canonicalizedQueryString.WriteByte('&')
		}
		canonicalizedQueryString.WriteString(c.percentEncode(k))
		canonicalizedQueryString.WriteByte('=')
		canonicalizedQueryString.WriteString(c.percentEncode(params[k]))
	}

	// 3. 构建待签名字符串
	stringToSign := "GET&" + c.percentEncode("/") + "&" + c.percentEncode(canonicalizedQueryString.String())

	// 4. 计算签名
	key := []byte(c.config.AccessKeySecret + "&")
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return signature
}

func (c *Client) percentEncode(s string) string {
	encoded := url.QueryEscape(s)
	encoded = strings.ReplaceAll(encoded, "+", "%20")
	encoded = strings.ReplaceAll(encoded, "*", "%2A")
	encoded = strings.ReplaceAll(encoded, "%7E", "~")
	return encoded
}

func (c *Client) buildQuery(params map[string]string) string {
	var pairs []string
	for k, v := range params {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
	}
	return strings.Join(pairs, "&")
}

func (c *Client) validateRequest(req *Request) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	if len(req.PhoneNumbers) == 0 {
		return errors.New("phone numbers cannot be empty")
	}

	if req.SignName == "" {
		return errors.New("sign name cannot be empty")
	}

	if req.TemplateId == "" {
		return errors.New("template id cannot be empty")
	}

	return nil
}
