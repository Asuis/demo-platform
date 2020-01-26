package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"demo-plaform/services/user"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Auth")
		claims, err := user.ParseToken(tokenString)
		if err!=nil {
			ctx.Status(http.StatusNonAuthoritativeInfo)
			return
		}
		ctx.Set("u", claims)
		//get auth
		ctx.Next()
	}
}