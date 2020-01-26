package db

import (
	"time"
)

type Mirror struct {
	ID          int64
	RepoID      int64
	Repo        *Repository `xorm:"-" json:"-"`
	Interval    int         // Hour.
	EnablePrune bool        `xorm:"NOT NULL DEFAULT true"`

	// Last and next sync time of Git data from upstream
	LastSync     time.Time `xorm:"-" json:"-"`
	LastSyncUnix int64     `xorm:"updated_unix"`
	NextSync     time.Time `xorm:"-" json:"-"`
	NextSyncUnix int64     `xorm:"next_update_unix"`

	address string `xorm:"-" json:"-"`
}
