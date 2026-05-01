package middleware

import (
	"context"
	"net/http"

	"moremei/ai-saas/cmd/tenant/api/internal/config"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/rest/token"
)

const (
	jwtAudience  = "aud"
	jwtExpire    = "exp"
	jwtId        = "jti"
	jwtIssueAt   = "iat"
	jwtIssuer    = "iss"
	jwtNotBefore = "nbf"
	jwtSubject   = "sub"
)

// TokenCheckMiddleware 检查用户登录状态, 未登录也会触发 api logic
type TokenCheckMiddleware struct {
	Config *config.Config
}

func NewTokenCheckMiddleware(conf *config.Config) *TokenCheckMiddleware {
	return &TokenCheckMiddleware{
		Config: conf,
	}
}

func (m *TokenCheckMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	parser := token.NewTokenParser()

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tok, err := parser.ParseToken(r, m.Config.Auth.AccessSecret, "")
		if err != nil {
			next(w, r.WithContext(ctx))
			return
		}
		if !tok.Valid {
			next(w, r.WithContext(ctx))
			return
		}

		claims, ok := tok.Claims.(jwt.MapClaims)
		if !ok {
			next(w, r.WithContext(ctx))
			return
		}

		for k, v := range claims {
			switch k {
			case jwtAudience, jwtExpire, jwtId, jwtIssueAt, jwtIssuer, jwtNotBefore, jwtSubject:
				// ignore the standard claims
			default:
				ctx = context.WithValue(ctx, k, v)
			}
		}

		next(w, r.WithContext(ctx))
	}
}
