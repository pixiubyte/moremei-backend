package xxl_job

type Config struct {
	ServerAddr   string `json:"ServerAddr"`
	AccessToken  string `json:"AccessToken,default=default_token"`
	RegistryKey  string `json:"RegistryKey"`
	ExecutorPort string `json:"ExecutorPort,default=9999"`
	ExecutorIp   string `json:"ExecutorIp,optional"`
	Timeout      int64  `json:"Timeout,default=5"`
}
