package proxy

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
)


func ReverseProxy() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		//todo check target
		target := ctx.Param("action")

		director := func(req *http.Request) {
			r := ctx.Request
			req = r
			req.URL.Scheme = "http"
			req.URL.Host = target
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	}
}