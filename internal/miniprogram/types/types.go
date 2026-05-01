package types

type CommonResponse struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type GetOpenidResponse struct {
	CommonResponse
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
}

type AccessTokenResponse struct {
	CommonResponse
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type GetMobileResponse struct {
	CommonResponse
	PhoneInfo struct {
		PhoneNumber     string `json:"phoneNumber"`
		PurePhoneNumber string `json:"purePhoneNumber"`
		CountryCode     string `json:"countryCode"`
		Watermark       struct {
			Timestamp int    `json:"timestamp"`
			Appid     string `json:"appid"`
		} `json:"watermark"`
	} `json:"phone_info"`
}

type SendMessageRequest struct {
	Touser           string      `json:"touser"`
	TemplateId       string      `json:"template_id"`
	Page             string      `json:"page"`
	Data             interface{} `json:"data"`
	MiniprogramState string      `json:"miniprogram_state"`
	Lang             string      `json:"lang"`
}

type GetUnlimitedQRCodeRequest struct {
	Scene      string `json:"scene"`
	Width      int64  `json:"width,optional,default=430"`
	Hyaline    bool   `json:"is_hyaline,optional,default=true"`
	Page       string `json:"page,optional,default=pages/home/home"`
	CheckPath  bool   `json:"check_path,optional,default=true"`
	EnvVersion string `json:"env_version,optional,default=release"`
}

type StableToken struct {
	GrantType    string `json:"grant_type"`
	Appid        string `json:"appid"`
	Secret       string `json:"secret"`
	ForceRefresh bool   `json:"force_refresh"`
}

type GetUrlLinkRequest struct {
	Path           string `json:"path,omitempty"`
	Query          string `json:"query,omitempty"`
	ExpireType     int64  `json:"expire_type,omitempty"`     // 默认值0.小程序 URL Link 失效类型，失效时间：0，失效间隔天数：1
	ExpireTime     int64  `json:"expire_time,omitempty"`     // 到期失效的 URL Link 的失效时间，为 Unix 时间戳。生成的到期失效 URL Link 在该时间前有效。最长有效期为30天。expire_type 为 0 必填
	ExpireInterval int64  `json:"expire_interval,omitempty"` // 到期失效的 URL Link 的失效间隔天数。生成的到期失效 URL Link 在该时间间隔到达前有效。最长有效期为30天。expire_type 为 1 必填
	EnvVersion     string `json:"env_version,omitempty"`     // 默认值"release"。要打开的小程序版本。正式版为 "release"，体验版为"trial"，开发版为"develop"，仅在微信外打开时生效。
}

type GetUrlLinkResponse struct {
	CommonResponse
	UrlLink string `json:"url_link"` // 生成的 URL Link
}
