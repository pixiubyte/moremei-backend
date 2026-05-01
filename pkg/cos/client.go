package cos

import (
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	cc  *cos.Client
	sc  *sts.Client
	cfg *Config
}

func NewClient(cfg Config) *Client {
	var (
		bUrl *url.URL
		sUrl *url.URL
	)
	if cfg.BucketURL != "" {
		bUrl, _ = url.Parse(cfg.BucketURL)
		sUrl, _ = url.Parse(fmt.Sprintf("https://cos.%s.myqcloud.com", cfg.Region))
	} else {
		bUrl, _ = url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", cfg.Bucket, cfg.Region))
		sUrl, _ = url.Parse(fmt.Sprintf("https://cos.%s.myqcloud.com", cfg.Region))
	}

	hc := http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cfg.SecretID,
			SecretKey: cfg.SecretKey,
		},
	}

	return &Client{
		cc:  cos.NewClient(&cos.BaseURL{BucketURL: bUrl, ServiceURL: sUrl}, &hc),
		sc:  sts.NewClient(cfg.SecretID, cfg.SecretKey, nil),
		cfg: &cfg,
	}
}

func (c *Client) GetCosClient() *cos.Client {
	return c.cc
}

func (c *Client) GetStsClient() *sts.Client {
	return c.sc
}

func (c *Client) GetUploadCredential(filePaths []string, expired time.Duration) (*sts.CredentialResult, error) {
	resource := make([]string, len(filePaths))
	for i, filePath := range filePaths {
		resource[i] = fmt.Sprintf(
			"qcs::cos:%s:uid/%s:%s/%s",
			c.cfg.Region,
			c.cfg.AppId,
			c.cfg.Bucket,
			filePath,
		)
	}
	opt := sts.CredentialOptions{
		Policy: &sts.CredentialPolicy{
			Statement: []sts.CredentialPolicyStatement{
				{
					Action: []string{
						"name/cos:PostObject",
						"name/cos:PutObject",
					},
					Effect:   "allow",
					Resource: resource,
				},
			},
		},
		Region:          c.cfg.Region,
		DurationSeconds: int64(expired.Seconds()),
	}

	return c.sc.GetCredential(&opt)
}
