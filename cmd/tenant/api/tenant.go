package main

import (
	"flag"
	"fmt"

	"moremei/ai-saas/cmd/tenant/api/internal/config"
	"moremei/ai-saas/cmd/tenant/api/internal/handler"
	"moremei/ai-saas/cmd/tenant/api/internal/middleware"
	"moremei/ai-saas/cmd/tenant/api/internal/svc"
	"moremei/ai-saas/pkg/x/http"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/tenant-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf, rest.WithCors("*"))
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	// 全局中间件
	// 全局中间件, 判断注册用户状态
	server.Use(middleware.NewUserCheckMiddleware(ctx.UserModel).Handle)
	// 注册路由
	handler.RegisterHandlers(server, ctx)
	// 设置自定义成功拦截器
	httpx.SetOkHandler(http.OkHandle)
	// 设置自定义错误拦截器
	httpx.SetErrorHandler(http.ErrorHandle)
	// 设置自定义请求验证器
	httpx.SetValidator(http.NewValidator())
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
