package utils

import (
	"context"
	"encoding/json"

	"moremei/ai-saas/cmd/tenant/api/internal/consts"
	errorx "moremei/ai-saas/pkg/x/error"

	"github.com/zeromicro/go-zero/core/logx"
)

type MeiAngel struct {
	UserId  uint64
	AngelId string
}

type AuthUser struct {
	UserId   uint64
	MeiAngel *MeiAngel
}

func GetAuthUserCtx(ctx context.Context) (*AuthUser, error) {
	var (
		meiAngel *MeiAngel
		userId   uint64
	)
	val := ctx.Value(consts.UserIDCtxKey)
	if val == nil {
		return nil, errorx.UnauthorizedError
	}
	switch value := val.(type) {
	case float64:
		userId = uint64(value)
	case json.Number:
		iv, err := value.Int64()
		if err != nil {
			logx.Errorv(err)
			return nil, errorx.UnauthorizedError
		}
		userId = uint64(iv)
	default:
		logx.Errorf("%s类型未知：%v", consts.UserIDCtxKey, val)
		return nil, errorx.UnauthorizedError
	}
	if val := ctx.Value(consts.MeiAngelCtxKey); val != nil {
		if ma, ok := val.(*MeiAngel); ok {
			if ma.UserId != 0 && ma.AngelId != "" {
				meiAngel = ma
			}
		}
	}

	// 校验用户ID与酶好天使用户ID是否一致
	if meiAngel != nil && meiAngel.UserId != userId {
		meiAngel = nil
	}

	return &AuthUser{
		UserId:   userId,
		MeiAngel: meiAngel,
	}, nil
}
