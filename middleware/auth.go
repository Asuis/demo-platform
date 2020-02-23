package middleware

import (
	"demo-platform/services/user"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		claims, err := user.ParseToken(tokenString)
		if err!=nil {
			ctx.Status(http.StatusUnauthorized)
			return
		}
		ctx.Set("u", claims)
		//get auth
		ctx.Next()
	}
}

func GitAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHead := ctx.GetHeader("Authorization")
		if len(authHead) == 0 {
			ctx.Status(http.StatusUnauthorized)
			return
		}
		auths := strings.Fields(authHead)
		if len(auths) != 2 || auths[0] != "Basic" {
			ctx.Status(http.StatusUnauthorized)
			return
		}
		authUsername, authPassword, err := BasicAuthDecode(auths[1])
		if err != nil {
			ctx.Status(http.StatusUnauthorized)
			return
		}
		_, err = user.Login(&user.LoginForm{
			Account: authUsername,
			Passwd:  authPassword,
		})
		if err != nil {
			ctx.Status(http.StatusUnauthorized)
			return
		}
	}
}

func BasicAuthDecode(encoded string) (string, string, error) {
	s, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", "", err
	}

	auth := strings.SplitN(string(s), ":", 2)
	return auth[0], auth[1], nil
}