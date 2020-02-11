package repo

import (
	"demo-platform/model/db"
	"errors"
	"fmt"
	"github.com/gogs/git-module"
	"os"
	"path"
	"strings"
	"time"
)

const BaseDir  = "/var/srv/git/"

type RepositoryInit struct {
	Name string `form:"name" json:"name" binding:"required"`
	Description string `form:"description" json:"description" binding:"required"`
	IsPrivate bool `form:"isPrivate" json:"isPrivate" binding:"required"`
	UseCustomAvatar bool `form:"useCustomAvatar" json:"useCustomAvatar" binding:"required"`
}

func Create(form *RepositoryInit, user *db.User) error  {

	session := db.Engine.NewSession()
	defer session.Close()


	p := path.Join(BaseDir, form.Name)

	has, err := session.Get(&db.Repository{Name:form.Name})

	if err != nil {
		return err
	}

	if has {
		return errors.New("the repository already exist")
	}

	now := time.Now()

	_ = session.Begin()

	_, err = session.Insert(&db.Repository{
		OwnerID:               user.Id,
		Owner:                 user,
		LowerName:             "",
		Name:                  form.Name,
		Description:           form.Description,
		Website:               "",
		DefaultBranch:         "master",
		Size:                  0,
		UseCustomAvatar:       form.UseCustomAvatar,
		NumWatches:            0,
		NumStars:              0,
		NumForks:              0,
		NumIssues:             0,
		NumClosedIssues:       0,
		NumOpenIssues:         0,
		NumPulls:              0,
		NumClosedPulls:        0,
		NumOpenPulls:          0,
		NumMilestones:         0,
		NumClosedMilestones:   0,
		NumOpenMilestones:     0,
		NumTags:               0,
		IsPrivate:             form.IsPrivate,
		IsBare:                true,
		IsMirror:              false,
		Mirror:                nil,
		EnableWiki:            false,
		AllowPublicWiki:       false,
		EnableExternalWiki:    false,
		ExternalWikiURL:       "",
		EnableIssues:          false,
		AllowPublicIssues:     false,
		EnableExternalTracker: false,
		ExternalTrackerURL:    "",
		ExternalTrackerFormat: "",
		ExternalTrackerStyle:  "",
		ExternalMetas:         nil,
		EnablePulls:           false,
		PullsIgnoreWhitespace: false,
		PullsAllowRebase:      false,
		IsFork:                false,
		Created:               now,
		CreatedUnix:           now.Unix(),
		Updated:               now,
		UpdatedUnix:           now.Unix(),
	})

	if strings.HasSuffix(p, ".git") {
		p += ".git"
	}

	err = git.InitRepository(p, true)

	if err != nil {
		_ = session.Rollback()
		return err
	}
	_ = session.Commit()

	return err
}


func Info(name string, user *db.User) (*db.Repository, error){

	repo := db.Repository{Name:name, OwnerID:user.Id}
	has, err := db.GetRepository(&repo)
	if err!=nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("the repository is not exist")
	}
	return &repo, nil
}

func List(user *db.User, page int, pageSize int, order string) (*[]db.Repository, error)  {
	repo := db.Repository{OwnerID:user.Id}
	list, err := db.ListRepository(&repo, page, pageSize, order)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func Clone(from string, to string) {}

func Fork() {}

func Del(name string, user *db.User) error {

	session := db.Engine.NewSession()

	defer session.Close()

	err := session.Begin()

	_, err = session.Delete(&db.Repository{
		Name: name,
		OwnerID:user.Id,
	})

	if err != nil {
		_ = session.Rollback()
		return fmt.Errorf("remove repo data failed")
	}
	if strings.HasSuffix(name, ".git") {
		name += ".git"
	}
	p := path.Join(BaseDir, name)
	err = os.RemoveAll(p)
	if err != nil {
		_ = session.Rollback()
		return fmt.Errorf("remove git file failed")
	}
	err = session.Commit()
	if err != nil {
		return err
	}
	return nil
}

func Download(path string) {}