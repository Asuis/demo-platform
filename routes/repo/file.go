package repo

import (
	"demo-platform/model/db"
	"demo-platform/services/repo"
	"demo-platform/services/user"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gogs/git-module"
	"io"
	"net/http"
	"strings"
)

type GitEntry struct {
	ID [20]byte
	Type string
	Name string
	IsDir bool
	IsSubModule bool
	IsLink bool
	Size int64
}

func SearchDir(ctx *gin.Context) {

	var has bool
	var err error
	sign, _ := ctx.Get("u")

	username := ctx.Param("username")
	r := ctx.Param("repo")
	relpath := ctx.Param("relpath")
	if relpath == "" {
		relpath = "/"
	}

	name := fmt.Sprintf("%s/%s", username, r)

	has, err = db.GetRepository(&db.Repository{
		OwnerID: sign.(*user.SignedData).Ac,
		Name:    name,
	})

	if err != nil || !has {
		ctx.Status(http.StatusNotFound)
		return
	}

	var data *git.Entries

	data, err = repo.SearchDir(name, relpath)

	if err != nil {

		ctx.JSON(http.StatusInternalServerError, err)
		return

	}
	var res []GitEntry
	for _,item := range *data {
		file := GitEntry{
			ID:   item.ID,
			Type: string(item.Type),
			Name: item.Name(),
			IsDir:item.IsDir(),
			IsSubModule:item.IsSubModule(),
			IsLink:item.IsLink(),
			Size:item.Size(),
		}
		res = append(res, file)
	}
	ctx.JSON(http.StatusOK, &res)
	
	return
}

type ContentType map[string] string


func GetRawFile(ctx *gin.Context) {
	var has bool
	var err error
	sign, _ := ctx.Get("u")

	username := ctx.Param("username")
	r := ctx.Param("repo")

	relpath := ctx.Param("relpath")

	i := strings.LastIndex(relpath, ".")
	if i != -1 {
		suffix := relpath[i:]
		if suffix =="png" {
			ctx.Header("Content-Type", "image/png")
		} else {
			ctx.Header("Content-Type", "text/plain")
		}
	}

	name := fmt.Sprintf("%s/%s", username, r)
	has, err = db.GetRepository(&db.Repository{
		OwnerID: sign.(*user.SignedData).Ac,
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
	reader, err := blob.Data()
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	_, _ = io.Copy(ctx.Writer, reader)
	return
}