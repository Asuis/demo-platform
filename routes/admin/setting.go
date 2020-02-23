package admin

import (
	"demo-platform/services/admin"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Users(ctx *gin.Context) {

	page := ctx.GetInt("page")
	pageSize := ctx.GetInt("pageSize")
	order := ctx.Param("order")

	data, err := admin.ListUser(page, pageSize, order)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, &data)
	return
}

func User(ctx *gin.Context) {
}

func Containers(ctx *gin.Context)  {
	page := ctx.GetInt("page")
	pageSize := ctx.GetInt("pageSize")
	order := ctx.Param("order")

	data, err := admin.ListContainer(page, pageSize, order)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, &data)
	return
}

func Container(ctx *gin.Context) {}

func Images(ctx *gin.Context) {
	page := ctx.GetInt("page")
	pageSize := ctx.GetInt("pageSize")
	order := ctx.Param("order")

	data, err := admin.ListRepository(page, pageSize, order)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, &data)
	return
}

func Repositories(ctx *gin.Context) {
	page := ctx.GetInt("page")
	pageSize := ctx.GetInt("pageSize")
	order := ctx.Param("order")

	data, err := admin.ListRepository(page, pageSize, order)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, &data)
	return
}

func Repository(ctx *gin.Context) {}

func Proxies(ctx *gin.Context) {}

func Proxy(ctx *gin.Context) {}