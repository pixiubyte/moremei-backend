package common

import (
	"net/http"

	"moremei/ai-saas/cmd/tenant/api/internal/logic/common"
	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	"moremei/ai-saas/pkg/x/error"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetAppLinkHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetAppLinkRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, error.WithCause(error.InvalidParamsError, err.Error()))
			return
		}

		l := common.NewGetAppLinkLogic(r.Context(), svcCtx)
		resp, err := l.GetAppLink(&req)

		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			if req.Redirect { // 重定向
				http.Redirect(w, r, resp.AppLink, http.StatusFound)
				return
			}
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
