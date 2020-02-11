package middleware

import (
	"demo-platform/services/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
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