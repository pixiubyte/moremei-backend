package http

import (
	"context"
	"errors"
	"net/http"
	"runtime/debug"
	"time"

	errorx "moremei/ai-saas/pkg/x/error"

	"github.com/zeromicro/go-zero/core/logx"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Time    int64  `json:"time"`
}

func OkHandle(ctx context.Context, data any) any {
	customErr := convert(nil)
	return &Response{
		Code:    customErr.GetCode(),
		Message: customErr.Error(),
		Data:    data,
		Time:    time.Now().Unix(),
	}
}

func ErrorHandle(err error) (int, any) {
	customErr := convert(err)
	if customErr != nil {
		if errors.Is(customErr, errorx.UnauthorizedError) {
			return http.StatusUnauthorized, errorx.UnauthorizedError
		} else if errors.Is(customErr, errorx.ForbiddenError) {
			return http.StatusForbidden, errorx.ForbiddenError
		}
		return http.StatusOK, &Response{
			Code:    customErr.GetCode(),
			Message: customErr.Error(),
			Data:    nil,
			Time:    time.Now().Unix(),
		}
	}

	return http.StatusInternalServerError, errorx.InternalServerError
}

func convert(err error) *errorx.Error {
	if err == nil {
		return errorx.NonError
	}

	var customErr *errorx.Error
	switch {
	case errors.As(err, &customErr):
		return customErr
	default:
		logx.Errorf("Unable to convert undefined error：%v, stack: %s", err, debug.Stack())
		return nil
	}
}
