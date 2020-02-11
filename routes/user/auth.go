package user

import (
	"demo-platform/services/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SignIn(ctx *gin.Context) {
	var json user.Login
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := user.SignIn(&json)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"t": token,
	})
	return
}

func SignUp(ctx *gin.Context) {
	var json user.Register
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := user.SignUp(&json)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"t": token,
	})
	return
}
