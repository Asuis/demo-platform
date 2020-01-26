package db

import "time"

type DockerImage struct {

	ID int64
	Sha1 string

	UserID int64
	Name string
	Desc string

	Author int64
	RepoPath string

	Created     time.Time `xorm:"-" json:"-"`
	CreatedUnix int64
	Updated     time.Time `xorm:"-" json:"-"`
	UpdatedUnix int64

}
