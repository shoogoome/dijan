package view

import (
	"dijan/core/view/auth"
	"github.com/kataras/iris"
)

func DijanViewBase(ctx iris.Context) auth.DijanAuthAuthorization {
	return auth.NewAuthAuthorization(&ctx)
}


