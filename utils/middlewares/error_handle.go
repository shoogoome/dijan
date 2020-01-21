package middlewares

import (
	"dijan/models"
	"fmt"
	"github.com/kataras/iris"
)

func AbnormalHandle(ctx iris.Context) {
	defer func() {
		re := recover()
		if re == nil {
			return
		}
		ctx.StatusCode(iris.StatusInternalServerError)
		switch result := re.(type) {
		case models.RestfulAPIResult:
			ctx.JSON(result)
		default:
			panic(re)
			ctx.JSON(models.RestfulAPIResult{
				Status: false,
				ErrCode: 500,
				Message: fmt.Sprintf("%v", re),
			})
		}
	}()
	ctx.Next()
}
