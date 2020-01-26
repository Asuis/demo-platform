package db

type Collaboration struct {
	ID     int64
	RepoID int64      `xorm:"UNIQUE(s) INDEX NOT NULL"`
	UserID int64      `xorm:"UNIQUE(s) INDEX NOT NULL"`
	Mode   AccessMode `xorm:"DEFAULT 2 NOT NULL"`
}
