package context

import (
	"moremei/ai-saas/internal/wecom/config"
	"moremei/ai-saas/internal/wecom/credential"
)

// Context struct
type Context struct {
	*config.Config
	credential.AccessTokenHandle
}
