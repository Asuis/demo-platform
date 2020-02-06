package db

import "time"

type DockerStatus int

const (
	DockerRunning = iota //0
	DockerStopping //1
	DockerStoped //2
	DockerBuilding //3
	DockerBuildSuccess //4
	DockerBuildFail //5
	)

type  DockerContainer struct {
	ID int64

	ContainerID string
	OwnerID int64
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