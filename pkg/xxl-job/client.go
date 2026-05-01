package xxl_job

import (
	"context"

	"moremei/ai-saas/third_party/xxl-job-executor-go"

	"github.com/zeromicro/go-zero/core/logx"
)

type Executor struct {
	exec xxl.Executor
	cfg  *Config
}

type RegisterFunc func() map[string]xxl.TaskFunc

func NewExecutor(cfg Config) *Executor {
	exec := xxl.NewExecutor(
		xxl.ServerAddr(cfg.ServerAddr),
		xxl.AccessToken(cfg.AccessToken), //请求令牌(默认为空)
		xxl.RegistryKey(cfg.RegistryKey), //执行器名称
		xxl.SetTimeout(cfg.Timeout),
		xxl.SetLogger(NewLogger(logx.WithContext(context.Background()))),
		xxl.ExecutorPort(cfg.ExecutorPort), //默认9999（非必填）
		xxl.ExecutorIp(cfg.ExecutorIp),
	)
	exec.Use(customMiddleware)
	exec.Init()

	return &Executor{
		exec: exec,
		cfg:  &cfg,
	}
}

func (e *Executor) RegisterTasks(registerFunc RegisterFunc) {
	for k, f := range registerFunc() {
		e.exec.RegTask(k, f)
	}
}

func (e *Executor) Start() {
	e.exec.Start()
}

func (e *Executor) Stop() {
	e.exec.Stop()
}

// 自定义中间件
func customMiddleware(tf xxl.TaskFunc) xxl.TaskFunc {
	return func(cxt context.Context, param *xxl.RunReq) string {
		logx.Info("xxljob middleware start")
		res := tf(cxt, param)
		logx.Info("xxljob middleware end")
		return res
	}
}
