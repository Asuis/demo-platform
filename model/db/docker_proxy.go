package db

import "time"

type DockerProxy struct {
	ID int64
	ContainerID int64
	ContainerName string
	ServerID string
	Port uint16
	Host string

	Created     time.Time `xorm:"-" json:"-"`
	CreatedUnix int64
	Updated     time.Time `xorm:"-" json:"-"`
	UpdatedUnix int64
}