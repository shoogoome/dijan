package ui

import (
	"dijan/core/view/auth"
	"dijan/utils/templates"
	"github.com/flosch/pongo2"
	"github.com/kataras/iris"
)

// 首页
func ManagerHome(ctx iris.Context, auth auth.DijanAuthAuthorization) {
	if auth.IsLogin() {
		ctx.Redirect("/overview")
	}
	ctx.HTML(templatesUtils.Render("index.html", pongo2.Context{
		"verification": false,
	}))
}
