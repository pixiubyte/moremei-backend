package config

type Auth struct {
	AccessSecret string
	AccessExpire int64
}

// 彩云天气api
type CaiYunApp struct {
	Url   string
	Token string
}

// 腾讯位置服务
type QQLbs struct {
	Url string `json:",default=https://apis.map.qq.com/ws"`
	Key string
}

type TestAccount struct {
	Mobile string
	Code   string
}

// 腾讯云语音识别
type TencentAsr struct {
	Url       string
	AppId     string
	SecretId  string
	SecretKey string
}
