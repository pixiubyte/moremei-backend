package captcha

type Config struct {
	SecretId         string `json:","`
	SecretKey        string `json:","`
	Region           string `json:",default=ap-shanghai"`
	CaptchaAppId     uint64 `json:","`
	CaptchaAppSecret string `json:","`
	Expire           int64  `json:",default=300"` // 单位秒
}
