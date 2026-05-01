package sms

type Config struct {
	SecretId  string `json:"SecretId"`
	SecretKey string `json:"SecretKey"`
	AppId     string `json:"AppId"`
	Endpoint  string `json:"Endpoint,default=sms.tencentcloudapi.com"`
	Region    string `json:"Region,default=ap-nanjing,options=ap-beijing|ap-guangzhou|ap-nanjing"`
}

type Template struct {
	AppId    string `json:"AppId"`
	SignName string `json:"SignName"`
	TplId    string `json:"TplId"`
}

type Templates map[string]Template
