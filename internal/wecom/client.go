package wecom

import (
	"moremei/ai-saas/internal/wecom/addresslist"
	"moremei/ai-saas/internal/wecom/app"
	"moremei/ai-saas/internal/wecom/cache"
	"moremei/ai-saas/internal/wecom/config"
	"moremei/ai-saas/internal/wecom/context"
	"moremei/ai-saas/internal/wecom/credential"
	"moremei/ai-saas/internal/wecom/externalcontact"
)

type Client struct {
	cfg            *config.Config
	ecClient       *externalcontact.Client
	oaClient       *app.Client
	cache          cache.Cache
	cacheKeyPrefix string
	ctx            *context.Context
}

func NewWecom(cfg *config.Config, cache cache.Cache, cacheKeyPrefix string) *Client {
	var oaClient *app.Client
	if cfg.OauthApp != nil {
		oaClient = app.NewClient(&context.Context{
			Config:            cfg,
			AccessTokenHandle: credential.NewWorkAccessToken(cfg.CorpID, cfg.OauthApp.Secret, cacheKeyPrefix, cache),
		})
	}

	ctx := &context.Context{
		Config:            cfg,
		AccessTokenHandle: credential.NewWorkAccessToken(cfg.CorpID, cfg.ExternalContact.Secret, cacheKeyPrefix, cache),
	}

	return &Client{
		cfg:            cfg,
		ecClient:       externalcontact.NewClient(ctx),
		oaClient:       oaClient,
		cache:          cache,
		cacheKeyPrefix: cacheKeyPrefix,
		ctx:            ctx,
	}
}

func (c *Client) GetExternalContact() *externalcontact.Client {
	return c.ecClient
}

func (c *Client) GetOauthApp() *app.Client {
	return c.oaClient
}

func (c *Client) GetAddressList() *addresslist.Client {
	return addresslist.NewClient(c.ctx)
}
