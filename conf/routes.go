package conf

import (
	"demo-platform/middleware"
	"demo-platform/routes/docker"
	"demo-platform/routes/proxy"
	repo2 "demo-platform/routes/repo"
	"demo-platform/routes/user"
	"demo-platform/services/repo"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter() *gin.Engine {

	r := gin.Default()

	v1 := r.Group("/v1")
	{
		gitRoutes := v1.Group("/git/:username/:repo")
		{
			gitRoutes.
				OPTIONS("/", func(ctx *gin.Context) {
					ctx.Status(http.StatusOK)
					return
				}).
				POST("/git-upload-pack", repo.ServiceUploadPack).
				POST("/git-receive-pack", repo.ServiceReceivePack).
				GET("/info/refs", repo.GetInfoRefs).
				GET("/HEAD", repo.GetTextFile).
				GET("/objects/*action", repo.GetObject)
		}

		userRoutes := v1.Group("/usr")
		{
			userRoutes.
				OPTIONS("/", func(ctx *gin.Context) {
				ctx.Status(http.StatusOK)
				return
			}).
				POST("/sign_in", user.SignIn).//ok
				POST("/sign_up", user.SignUp)//ok
		}

		repoRoutes := v1.Group("/repo")
		{
			repoRoutes.OPTIONS("/", func(ctx *gin.Context) {
				ctx.Status(http.StatusOK)
				return
			}).
				POST("/create", middleware.Auth(), repo2.CreateRepository).
				GET("/info/:username/:repo", middleware.Auth(), repo2.GetRepoInfo).
				DELETE("/:username/:repo", middleware.Auth(), repo2.Delete).
				PUT("/:username/:repo", middleware.Auth(), repo2.Setting).
				GET("/branch/:username/:repo", middleware.Auth(), repo2.GetRepoBranches).
				GET("/commit/:username/:repo", middleware.Auth()).
				GET("/list/:pageSize/:page/:order", middleware.Auth(), repo2.List).
				GET("/mylist/:pageSize/:page/:order", middleware.Auth(), repo2.List).
				GET("/file/:username/:repo/*relpath", middleware.Auth(), repo2.SearchDir).
				GET("/rawfile/:username/:repo/*relpath", middleware.Auth(), repo2.GetRawFile)
		}

		cloudRoutes := v1.Group("/cloud")
		{
			cloudRoutes.OPTIONS("/", func(ctx *gin.Context) {
				ctx.Status(http.StatusOK)
				return
			}).
				POST("/create", middleware.Auth(), docker.CreateDocker).
				DELETE("/:cloud_id",middleware.Auth(), docker.DelDocker).
				PUT("/:cloud_id", middleware.Auth(), docker.UpdateDocker).
				POST("/list", middleware.Auth(), docker.ListDocker).
				POST("/action/run",middleware.Auth(), docker.StartDocker).
				POST("/action/stop",middleware.Auth(), docker.StopDocker).
				POST("/action/restart").
				GET("/info/:cloud_id",middleware.Auth(), docker.InfoDocker).
				GET("/stat/:cloud_id",middleware.Auth(), docker.StatDocker).
				GET("/console/:cloud_id",middleware.Auth(), docker.AttachDocker)
		}

		proxyRoutes := v1.Group("/proxy")
		{
			proxyRoutes.Any("/*action", proxy.ReverseProxy())
		}

		adminRoutes := v1.Group("/admin")
		{
			adminRoutes.OPTIONS("/", func(ctx *gin.Context) {
				ctx.Status(http.StatusOK)
				return
			}).
				GET("/users").
				GET("/user/:id").
				GET("/containers").
				GET("/container/:id").
				GET("/images").
				GET("/repositories").
				GET("/repo/:username/:repo").
				GET("/auth", middleware.Auth(), func(ctx *gin.Context) {
				sign, _ := ctx.Get("u")
				fmt.Println(sign)
			})
		}
	}

	return r

}
