package service

import (
	"dijan/service/ui"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
)

func RegisterUIRouters(app *iris.Application) {

	// 后台路由
	app.Get("/", hero.Handler(ui.ManagerHome))
	app.Post("/login", hero.Handler(ui.Login))
	app.Get("/logout", hero.Handler(ui.Logout))
	app.Get("/overview", hero.Handler(ui.Overview))
	app.Get("/node/{node:string}", hero.Handler(ui.NodeUI))
	app.Get("/static/{fileType:string}/{fileName:string}", hero.Handler(ui.StaticResource))
	app.Get("/static/{fileType:string}/{fileName:string}", hero.Handler(ui.StaticResource))
}
