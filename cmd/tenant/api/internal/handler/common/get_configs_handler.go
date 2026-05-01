package common

import (
	"net/http"

	"moremei/ai-saas/cmd/tenant/api/internal/logic/common"
	"moremei/ai-saas/cmd/tenant/api/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetConfigsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := common.NewGetConfigsLogic(r.Context(), svcCtx)
		resp, err := l.GetConfigs()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
