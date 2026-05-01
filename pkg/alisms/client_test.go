package alisms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_Send(t *testing.T) {
	accessKeyId := ""
	accessKeySecret := ""

	client := MustNewClient(&Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		Region:          "cn-hangzhou",
	})

	tests := []struct {
		name    string
		req     *Request
		wantErr bool
	}{
		{
			name: "valid single phone",
			req: &Request{
				SignName:     "爱护健康",          // 替换为你的签名
				TemplateId:   "SMS_478475691", // 替换为你的模板ID
				PhoneNumbers: []string{"+8619921209365"},
				Params: map[string]interface{}{
					"code": "123456",
				},
			},
			wantErr: false,
		},
		{
			name: "valid multiple phones",
			req: &Request{
				SignName:     "测试签名",          // 替换为你的签名
				TemplateId:   "SMS_123456789", // 替换为你的模板ID
				PhoneNumbers: []string{"", ""},
				Params: map[string]interface{}{
					"code": "123456",
				},
			},
			wantErr: false,
		},
		{
			name: "empty phone numbers",
			req: &Request{
				SignName:     "测试签名",
				TemplateId:   "SMS_123456789",
				PhoneNumbers: []string{},
				Params: map[string]interface{}{
					"code": "123456",
				},
			},
			wantErr: true,
		},
		{
			name: "empty sign name",
			req: &Request{
				SignName:     "",
				TemplateId:   "SMS_478475691",
				PhoneNumbers: []string{"19921209365"},
				Params: map[string]interface{}{
					"code": "123456",
				},
			},
			wantErr: true,
		},
		{
			name: "empty template id",
			req: &Request{
				SignName:     "测试签名",
				TemplateId:   "",
				PhoneNumbers: []string{""},
				Params: map[string]interface{}{
					"code": "123456",
				},
			},
			wantErr: true,
		},
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.Send(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)

			// 检查每个手机号的响应
			for _, phone := range tt.req.PhoneNumbers {
				r, ok := resp[phone]
				assert.True(t, ok)
				assert.Equal(t, phone, r.PhoneNumber)
				assert.NotEmpty(t, r.RequestId)
				assert.NotEmpty(t, r.Code)
				assert.NotEmpty(t, r.Message)

				if r.Result {
					assert.Equal(t, SendSuccessCode, r.Code)
					assert.NotEmpty(t, r.BizId)
				}
			}
		})
	}
}

func TestClient_validateRequest(t *testing.T) {
	client := &Client{}

	tests := []struct {
		name    string
		req     *Request
		wantErr bool
	}{
		{
			name: "valid request",
			req: &Request{
				SignName:     "测试签名",
				TemplateId:   "SMS_123456789",
				PhoneNumbers: []string{"13800138000"},
				Params: map[string]interface{}{
					"code": "123456",
				},
			},
			wantErr: false,
		},
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
		},
		{
			name: "empty phone numbers",
			req: &Request{
				SignName:     "测试签名",
				TemplateId:   "SMS_123456789",
				PhoneNumbers: []string{},
			},
			wantErr: true,
		},
		{
			name: "empty sign name",
			req: &Request{
				SignName:     "",
				TemplateId:   "SMS_123456789",
				PhoneNumbers: []string{"13800138000"},
			},
			wantErr: true,
		},
		{
			name: "empty template id",
			req: &Request{
				SignName:     "测试签名",
				TemplateId:   "",
				PhoneNumbers: []string{"13800138000"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.validateRequest(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
