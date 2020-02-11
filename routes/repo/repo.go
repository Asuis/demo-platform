package repo

import (
	"demo-platform/model/db"
	"demo-platform/services/repo"
	"demo-platform/services/user"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CreateRepository(ctx *gin.Context) {

	var json repo.RepositoryInit

	sign, _ := ctx.Get("u")

	s := sign.(user.SignedData)

	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	usr := db.User{Id:s.Ac}

	_, err := db.GetByUser(&usr)

	err = repo.Create(&json, &usr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)

	return
}

func GetRepoInfo(ctx *gin.Context) {

	sign, _ := ctx.Get("u")

	username := ctx.Param("username")

	r := ctx.Param("repo")

	repository, err := repo.Info(username + "/"+r, &db.User{
		Id: sign.(user.SignedData).Ac,
	})

	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, repository)
	return
}

func List(ctx *gin.Context) {
	
	sign, _ := ctx.Get("u")

	pageSize, err := strconv.Atoi(ctx.Param("pageSize"))
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	page, err := strconv.Atoi(ctx.Param("page"))
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	order := ctx.Param("order")

	list, err := repo.List(&db.User{
		Id: sign.(user.SignedData).Ac,
	}, page, pageSize, order)

	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, list)

	return
}

func Setting(ctx *gin.Context) {
	return
}

func Delete(ctx *gin.Context) {

	sign, _ := ctx.Get("u")

	username := ctx.Param("username")
	
	r := ctx.Param("repo")

	err := repo.Del(username + "/" + r, &db.User{
		Id: sign.(user.SignedData).Ac,
	})

	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
	return
}
