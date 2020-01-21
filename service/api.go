package service

import (
	"dijan/service/api"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
)

func RegisterApiRouters(app *iris.Application) {

	apiRouter := app.Party("/api")

	apiRouter.Get("/node", hero.Handler(api.NodeStorageInfo))
	apiRouter.Get("/storage/get/{key:string}", hero.Handler(api.GetData))
	apiRouter.Post("/storage/set", hero.Handler(api.ModifyData))
	apiRouter.Delete("/storage/delete", hero.Handler(api.DeleteData))
	apiRouter.Get("/rebalance", hero.Handler(api.Rebalance))
	apiRouter.Get("/empty", hero.Handler(api.Empty))
	apiRouter.Get("/empty_search", hero.Handler(api.EmptySearch))
	apiRouter.Get("/empty_all", hero.Handler(api.EmptyAll))
	apiRouter.Get("/member", hero.Handler(api.GetMember))
}
