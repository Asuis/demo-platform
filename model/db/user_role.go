package db

import "time"

type Role struct {

	ID int64
	Name string
	Desc string
	Sort string

	Created time.Time `xorm:"-" json:"-"`
	CreatedUnix int64
	Updated time.Time `xorm:"-" json:"-"`
	UpdatedUnix int64
	Author int64

}

type UserRole struct {
	ID int64
	UserID int64
	Created time.Time `xorm:"-" json:"-"`
	CreatedUnix int64
	Updated time.Time `xorm:"-" json:"-"`
	UpdatedUnix int64
	Author int64
}