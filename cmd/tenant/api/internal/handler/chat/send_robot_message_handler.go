package chat

import (
	"net/http"

	"moremei/ai-saas/cmd/tenant/api/internal/logic/chat"
	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/cmd/tenant/api/internal/types"
	"moremei/ai-saas/pkg/x/error"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func SendRobotMessageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SendRobotMessageRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, error.WithCause(error.InvalidParamsError, err.Error()))
			return
		}

		l := chat.NewSendRobotMessageLogic(r.Context(), svcCtx, w)
		err := l.SendRobotMessage(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		}
	}
}
