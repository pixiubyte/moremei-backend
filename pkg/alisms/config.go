package alisms

type Config struct {
	AccessKeyId     string `json:"AccessKeyId"`
	AccessKeySecret string `json:"AccessKeySecret"`
	Region         string `json:"Region,default=cn-hangzhou"`
}

type Template struct {
	SignName string `json:"SignName"`
	TplId    string `json:"TplId"`
}

type Templates map[string]Template
