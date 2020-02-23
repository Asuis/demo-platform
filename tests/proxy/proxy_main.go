package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func t() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//todo check target
		//action := ctx.Param("action")
		origin, _ := url.Parse("http://127.0.0.1:8081")
		director := func(req *http.Request) {
			req.Header.Add("X-Forwarded-Host", req.Host)
			req.Header.Add("X-Origin-Host", origin.Host)
			req.URL.Scheme = "http"
			req.URL.Host = origin.Host
		}
		proxy := &httputil.ReverseProxy{Director: director}

		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

func main() {


	r := gin.Default()

	r.Any("/", t())

	_ = r.Run(":8000")

}