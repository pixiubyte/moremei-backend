package auth

import (
	"net/http"

	"moremei/ai-saas/cmd/tenant/api/internal/logic/auth"
	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	"moremei/ai-saas/pkg/x/error"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetTokenRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, error.WithCause(error.InvalidParamsError, err.Error()))
			return
		}

		l := auth.NewGetTokenLogic(r.Context(), svcCtx)
		resp, err := l.GetToken(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
