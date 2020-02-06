package repo

import (
	"demo-plaform/model/db"
	"demo-plaform/services/repo"
	"demo-plaform/services/user"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gogs/git-module"
	"net/http"
)

func SearchDir(ctx gin.Context) {

	var has bool
	var err error
	sign, _ := ctx.Get("u")

	username := ctx.Param("username")
	r := ctx.Param("repo")

	name := fmt.Sprintf("%s/%s", username, r)

	has, err = db.GetRepository(&db.Repository{
		OwnerID: sign.(user.SignedData).Ac,
		Name:    name,
	})

	if err != nil || !has {
		ctx.Status(http.StatusNotFound)
		return
	}

	var data *git.Entries

	data, err = repo.SearchDir(name)

	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, data)
	
	return
}

func GetRawFile(ctx gin.Context) {
	var has bool
	var err error
	sign, _ := ctx.Get("u")

	username := ctx.Param("username")
	r := ctx.Param("repo")

	relpath := ctx.Param("relpath")

	name := fmt.Sprintf("%s/%s", username, r)
	has, err = db.GetRepository(&db.Repository{
		OwnerID: sign.(user.SignedData).Ac,
		Name:    name,
	})

	if err != nil || !has {
		ctx.Status(http.StatusNotFound)
		return
	}
	blob, err := repo.GetRawFile(name, relpath)

	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, blob)

	return
}