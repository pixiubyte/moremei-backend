package middleware

import (
	"errors"
	"net/http"

	"moremei/ai-saas/cmd/tenant/api/internal/utils"
	"moremei/ai-saas/internal/model/tenant"
	"moremei/ai-saas/pkg/gormc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// UserCheckMiddleware 检查用户账号状态，拒绝禁用的用户
type UserCheckMiddleware struct {
	UserModel tenant.UserModel
}

func NewUserCheckMiddleware(userModel tenant.UserModel) *UserCheckMiddleware {
	return &UserCheckMiddleware{
		UserModel: userModel,
	}
}

func (m *UserCheckMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		authUser, err := utils.GetAuthUserCtx(r.Context())
		if err != nil { // 未登录接口，直接放行
			next(w, r)
			return
		}

		// 查询用户的信息
		user, err := m.UserModel.FindOne(r.Context(), authUser.UserId)
		if err != nil {
			if !errors.Is(err, gormc.ErrNotFound) {
				httpx.ErrorCtx(r.Context(), w, err)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if user.Status == 0 { // 禁用的用户
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
