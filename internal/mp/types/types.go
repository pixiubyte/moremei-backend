package types

type Template struct {
	Touser      string      `json:"touser"`
	TemplateId  string      `json:"template_id"`
	Url         string      `json:"url"`
	ClientMsgId string      `json:"client_msg_id"`
	Data        interface{} `json:"data"`
	Miniprogram Miniprogram `json:"miniprogram"`
}

type Miniprogram struct {
	Appid    string `json:"appid"`
	Pagepath string `json:"pagepath"`
}

type Common struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type OauthOpenidResp struct {
	Common
	AccessToken    string `json:"access_token"`
	ExpiresIn      int    `json:"expires_in"`
	RefreshToken   string `json:"refresh_token"`
	Openid         string `json:"openid"`
	Scope          string `json:"scope"`
	IsSnapshotuser int    `json:"is_snapshotuser"`
	Unionid        string `json:"unionid"`
}

type AccessTokenResp struct {
	Common
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type ComponentTokenRequest struct {
	ComponentAppid        string `json:"component_appid"`
	ComponentAppsecret    string `json:"component_appsecret"`
	ComponentVerifyTicket string `json:"component_verify_ticket"`
}

type ComponentTokenResponse struct {
	Common
	ComponentAccessToken string `json:"component_access_token"`
	ExpiresIn            int    `json:"expires_in"`
}

type PreAuthCodeRequest struct {
	ComponentAppid string `json:"component_appid"`
}

type PreAuthCodeResponse struct {
	Common
	PreAuthCode string `json:"pre_auth_code"`
	ExpiresIn   int    `json:"expires_in"`
}

type UserInfoResponse struct {
	Common
	Openid     string `json:"openid"`
	Nickname   string `json:"nickname"`
	Sex        int64  `json:"sex"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Country    string `json:"country"`
	Headimgurl string `json:"headimgurl"`
}

type UnionIDResponse struct {
	Common
	UnionId string `json:"unionid"`
}

type AccessTokenRequest struct {
	GrantType string `json:"grant_type"`
	Appid     string `json:"appid"`
	Secret    string `json:"secret"`
}
