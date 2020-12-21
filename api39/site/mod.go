package site

import (
	"github.com/kataras/iris/v12"
)

func Update(ctx iris.Context) {
	ctx.JSON(iris.Map{"message": "successful request"})
}
