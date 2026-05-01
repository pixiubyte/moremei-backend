package chat

import (
	"net/http"

	"moremei/ai-saas/cmd/tenant/api/internal/logic/chat"
	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	"moremei/ai-saas/pkg/x/error"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetAiRobotsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetAiRobotsRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, error.WithCause(error.InvalidParamsError, err.Error()))
			return
		}

		l := chat.NewGetAiRobotsLogic(r.Context(), svcCtx)
		resp, err := l.GetAiRobots(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
