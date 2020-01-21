package service

import (
	"dijan/core/view"
	"dijan/utils"
	"dijan/utils/middlewares"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
)


// 初始化http服务路由映射
func initRouter(app *iris.Application) {
	RegisterUIRouters(app)
	RegisterApiRouters(app)
}


// 初始化http服务
func InitHttpService() {
	go func() {
		app := iris.New()
		// 注册控制器
		app.UseGlobal(middlewares.AbnormalHandle)
		hero.Register(view.DijanViewBase)

		initRouter(app)
		app.Run(iris.Addr(utils.GlobalSystemConfig.Server.HttpListenPort), iris.WithoutServerError(iris.ErrServerClosed))
	}()
}
