package config

type AppConfig struct {
	AgentID        string `json:"AgentId"`
	Secret         string `json:"Secret"`
	Token          string `json:"Token,optional"`
	RasPrivateKey  string `json:"RasPrivateKey,optional"`
	EncodingAESKey string `json:"EncodingAESKey,optional"`
}

type Config struct {
	CorpID          string     `json:"CorpID"`
	ExternalContact AppConfig  `json:"ExternalContact"`
	OauthApp        *AppConfig `json:"OauthApp,optional"`
}
