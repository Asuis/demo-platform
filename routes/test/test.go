package test

import (
	"demo-platform/services/user"
	"github.com/gin-gonic/gin"
)

func test(ctx *gin.Context) {
	var sign user.SignedData
	var has bool
	sign, has = ctx.Get("u")
	return
}