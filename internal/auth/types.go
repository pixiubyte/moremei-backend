package auth

const (
	ProviderApple     = "apple"
	ProviderWechat    = "wechat"
	ProviderThirdChat = "third_chat"
)

// AppleConf登录配置
type AppleConf struct {
	KeyId       string `json:"key_id"`
	TeamId      string `json:"team_id"`
	ClientId    string `json:"client_id"`             // 应用ID
	PrivateKey  string `json:"private_key"`           // 私钥, p8文件去掉头尾，去掉换行，base64编码
	AppSecret   string `json:"app_secret,optional"`   // 应用私钥
	RedirectUrl string `json:"redirect_url,optional"` // 回调地址
}

// 微信登录配置
type WechatConf struct {
	AppId       string `json:"app_id"`
	AppSecret   string `json:"app_secret"`
	RedirectUrl string `json:"redirect_url,optional"` // 回调地址
}
