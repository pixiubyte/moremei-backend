package chat

import (
	"moremei/ai-saas/cmd/tenant/api/internal/logic/chat"
	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetConversationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := chat.NewGetConversationsLogic(r.Context(), svcCtx)
		resp, err := l.GetConversations()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
