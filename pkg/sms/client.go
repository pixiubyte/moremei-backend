package sms

import (
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	commonError "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"github.com/zeromicro/go-zero/core/logx"
)

type Client struct {
	config *Config
	sc     *sms.Client
}

type Request struct {
	AppId        string // 非必需
	SignName     string
	TemplateId   string
	PhoneNumbers []string
	Params       []string
}

type Response struct {
	Result      bool
	PhoneNumber string
	RequestId   string
	Code        string
	Message     string
	SerialNo    string
}

const SendSuccessCode = "Ok"

func MustNewClient(cfg *Config) *Client {
	credential := common.NewCredential(cfg.SecretId, cfg.SecretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = cfg.Endpoint

	client, err := sms.NewClient(credential, cfg.Region, cpf)
	if err != nil {
		logx.Errorv(err)
	}

	return &Client{config: cfg, sc: client}
}

func (c *Client) Send(req *Request) (map[string]Response, error) {
	if req.AppId == "" {
		req.AppId = c.config.AppId
	}
	err := c.validateRequest(req)
	if err != nil {
		return nil, err
	}

	response, err := c.sc.SendSms(req.newSendSmsRequest())
	if err != nil {
		var sdkErr *commonError.TencentCloudSDKError
		if errors.Is(err, sdkErr) {
			return lo.SliceToMap[string, string, Response](req.PhoneNumbers, func(item string) (string, Response) {
				return item, Response{
					PhoneNumber: item,
					RequestId:   sdkErr.GetRequestId(),
					Result:      false,
					Code:        sdkErr.GetCode(),
					Message:     sdkErr.GetMessage(),
				}
			}), nil
		}

		return lo.SliceToMap[string, string, Response](req.PhoneNumbers, func(item string) (string, Response) {
			return item, Response{
				PhoneNumber: item,
				Result:      false,
				Message:     err.Error(),
			}
		}), nil
	}

	return lo.SliceToMap[string, string, Response](req.PhoneNumbers, func(phoneNumber string) (string, Response) {
		sendStatus, exist := lo.Find[*sms.SendStatus](response.Response.SendStatusSet, func(item *sms.SendStatus) bool {
			return item.PhoneNumber != nil && phoneNumber == *item.PhoneNumber
		})

		if exist {
			return phoneNumber, Response{
				PhoneNumber: phoneNumber,
				RequestId:   *response.Response.RequestId,
				Result:      *sendStatus.Code == SendSuccessCode,
				Code:        *sendStatus.Code,
				Message:     *sendStatus.Message,
				SerialNo:    *sendStatus.SerialNo,
			}
		}

		return phoneNumber, Response{
			PhoneNumber: phoneNumber,
			RequestId:   *response.Response.RequestId,
			Result:      false,
			Message:     "Internal error: doesn't match any phone number",
		}
	}), nil
}

func (c *Client) validateRequest(msg *Request) error {
	if msg.AppId == "" {
		return errors.New("短信应用ID不能为空")
	}

	if msg.SignName == "" {
		return errors.New("短信签名不能为空")
	}

	if msg.TemplateId == "" {
		return errors.New("短信模板ID不能为空")
	}

	if len(msg.PhoneNumbers) == 0 {
		return errors.New("短信接收手机号不能为空")
	}

	return nil
}

func (req *Request) newSendSmsRequest() *sms.SendSmsRequest {
	request := sms.NewSendSmsRequest()
	request.SmsSdkAppId = common.StringPtr(req.AppId)
	request.SignName = common.StringPtr(req.SignName)
	request.TemplateId = common.StringPtr(req.TemplateId)
	request.TemplateParamSet = common.StringPtrs(req.Params)
	request.PhoneNumberSet = common.StringPtrs(req.PhoneNumbers)

	return request
}
