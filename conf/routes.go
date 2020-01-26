package conf

import (
	"demo-plaform/services/repo"
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
			userRoutes.OPTIONS("/", func(ctx *gin.Context) {
				ctx.Status(http.StatusOK)
				return
			}).POST("/sign_in")
		}
	}

	return r
}
