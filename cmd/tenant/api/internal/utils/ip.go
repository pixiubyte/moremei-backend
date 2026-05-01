package utils

import (
	"context"
	"moremei/ai-saas/cmd/tenant/api/internal/consts"
)

func GetRemoteIpCtx(ctx context.Context) string {
	tenant, ok := ctx.Value(consts.RemoteIpCtxKey).(string)
	if !ok {
		return ""
	}
	return tenant
}
