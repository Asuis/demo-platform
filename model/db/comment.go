package db

import "time"

type Comment struct {
	ID int64

	AuthorID int64
	Title string
	Content string
	SubComment *[] Comment

	BindTypeId int64



	Created     time.Time `xorm:"-" json:"-"`
	CreatedUnix int64
	Updated     time.Time `xorm:"-" json:"-"`
	UpdatedUnix int64
}

type CommentType struct {
	ID int64
	Desc string

}