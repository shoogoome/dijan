package ui

import (
	"dijan/core/view/auth"
	"dijan/utils"
	templatesUtils "dijan/utils/templates"
	"github.com/flosch/pongo2"
	"github.com/kataras/iris"
)

// 登录
func Login(ctx iris.Context, auth auth.DijanAuthAuthorization) {
	password := ctx.FormValue("password")
	if password == utils.GlobalSystemConfig.Server.SystemPassword {
		auth.SetCookie(true)
		ctx.Redirect("/overview")
	}
	ctx.HTML(templatesUtils.Render("index.html", pongo2.Context{
		"verification": true,
		"login": false,
	}))
}

// 登出
func Logout(ctx iris.Context, auth auth.DijanAuthAuthorization) {
	auth.SetCookie(false)
	ctx.Redirect("/")
}
