package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"moremei/ai-saas/cmd/tenant/api/internal/consts"
	"moremei/ai-saas/internal/model/tenant"
	"moremei/ai-saas/pkg/gormc"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

type TokenAuthMiddleware struct {
	tokenManageModel         tenant.TokenManageModel
	tokenRequestHistoryModel tenant.TokenRequestHistoryModel
	redis                    *redis.Redis
	name                     string
}

func NewTokenAuthMiddleware(name string, tokenManageModel tenant.TokenManageModel, tokenRequestHistoryModel tenant.TokenRequestHistoryModel, redis *redis.Redis) *TokenAuthMiddleware {
	return &TokenAuthMiddleware{
		tokenManageModel:         tokenManageModel,
		tokenRequestHistoryModel: tokenRequestHistoryModel,
		redis:                    redis,
		name:                     name,
	}
}

func (m *TokenAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token := r.Header.Get("token")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		key := fmt.Sprintf(consts.MiniProgramTokenManage, m.name, token)
		val, err := m.redis.Get(key)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if val == "-" {
			//不存token,3秒缓存
			http.Error(w, "token鉴权失败", http.StatusUnauthorized)
			return
		} else if val == "" {
			tokenManage, err := m.tokenManageModel.FindOneByToken(ctx, token)
			if err != nil {
				if errors.Is(err, gormc.ErrNotFound) {
					err = m.redis.Setex(key, "-", 3)
					if err != nil {
						http.Error(w, "token鉴权失败", http.StatusUnauthorized)
						return
					}
				}
				http.Error(w, "token鉴权失败", http.StatusUnauthorized)
				return
			} else {
				err = m.redis.Setex(key, tokenManage.Token, 300)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}
		} else if val != token {
			http.Error(w, "token鉴权失败", http.StatusUnauthorized)
			return
		}
		err = m.tokenRequestHistoryModel.Insert(ctx, &tenant.TokenRequestHistory{
			Token: token,
			Route: r.RequestURI,
		}, nil)
		if err != nil {
			logx.WithContext(ctx).Errorf("路由日志插入失败")
		}

		next(w, r)
	}
}
