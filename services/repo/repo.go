package repo

import (
	"demo-plaform/model/db"
	"github.com/gogs/git-module"
)

type RepositoryInit struct {
	Name string
	Description string
	IsPrivate bool
	UseCustomAvatar bool
}

func Create(path string, user *db.User) error  {
	err := git.InitRepository(path, true)
	if err != nil {
		return err
	}
	return err
}

func Clone(from string, to string) {

}

func Fork() {

}

func Del() {}

func Download(path string) {}