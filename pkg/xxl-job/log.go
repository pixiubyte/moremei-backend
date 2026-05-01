package xxl_job

import (
	"moremei/ai-saas/third_party/xxl-job-executor-go"

	"github.com/zeromicro/go-zero/core/logx"
)

type logger struct {
	l logx.Logger
}

func NewLogger(l logx.Logger) xxl.Logger {
	return &logger{l: l}
}

func (l *logger) Info(format string, args ...interface{}) {
	l.l.Info(format, args)
}

func (l *logger) Error(format string, args ...interface{}) {
	l.l.Error(format, args)
}
