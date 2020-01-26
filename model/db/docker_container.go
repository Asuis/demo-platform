package db

import "time"

type DockerStatus int

const (
	DockerRunning = iota
	DockerStopping
	DockerStoped
	DockerBuilding
	DockerBuildSuccess
	DockerBuildFail
	)

type  DockerContainer struct {
	ID int64
	UserID int64
	Name string
	Desc string

	ImageID int64
	ImageSha1 string

	Status DockerStatus
	Port []* int
	ServerID int64

	LogPath string

	Created     time.Time `xorm:"-" json:"-"`
	CreatedUnix int64
	Updated     time.Time `xorm:"-" json:"-"`
	UpdatedUnix int64
}