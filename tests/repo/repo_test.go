package repo

import (
	"demo-platform/model/db"
	"demo-platform/services/repo"
	"log"
	"testing"
)

func TestCreateRepo(t *testing.T) {

	err := db.SetupDatabase()


	var usr = db.User{
		Id: 1,
	}
	_, err = db.GetByUser(&usr)

	err = repo.Create(&repo.RepositoryInit{
		Name:            "asuis/test2",
		Description:     "test2",
		IsPrivate:       false,
		UseCustomAvatar: false,
	}, &usr)
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func TestGetRepoList(t *testing.T) {
	err := db.SetupDatabase()

	var usr = db.User{
		Id: 1,
	}
	list, err:= repo.List(&usr, 0, 10, "ID")
	if err != nil {
		log.Fatalf("%v", err)
	}

	for _, item := range *list {
		log.Println(item.Name)
	}
}