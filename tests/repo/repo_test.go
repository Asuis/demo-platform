package repo

import (
	"demo-plaform/model/db"
	"demo-plaform/services/repo"
	"log"
	"testing"
)

func TestCreateRepo(t *testing.T) {
	var usr = db.User{
		Id: 1,
	}
	_, err := db.GetByUser(&db.User{
		Id: 1,
	})

	err = repo.Create(&repo.RepositoryInit{
		Name:            "asuis/test2.git",
		Description:     "test2",
		IsPrivate:       false,
		UseCustomAvatar: false,
	}, &usr)
	if err != nil {
		log.Fatalf("%v", err)
	}
}

