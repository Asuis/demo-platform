package proxy

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func ReverseProxy() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		//todo check target
		//action := ctx.Param("action")
		u := ctx.Request.RequestURI
		u = strings.TrimPrefix(u, "/v1/proxy")
		ctx.Request.RequestURI = u
		ctx.Request.URL.Path = u
		origin, _ := url.Parse("http://127.0.0.1:8081")
		director := func(req *http.Request) {
			req.Header.Add("X-Forwarded-Host", req.Host)
			req.Header.Add("X-Origin-Host", origin.Host)
			req.URL.Scheme = "http"
			req.URL.Host = origin.Host
		}

		proxy := &httputil.ReverseProxy{Director: director, ModifyResponse: func(res *http.Response) error {
			return nil
		}}

		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	}
}