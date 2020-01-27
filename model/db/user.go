package db

import (
	"time"
)

type UserType int

const (
	UserTypeIndividual UserType = iota // Historic reason to make it starts at 0.
	UserTypeOrganization
)

type User struct {
	Id int64
	Name string
	Passwd string `xorm:"varchar(200)"`

	// Email is the primary email address (to be used for communication)
	Email       string `xorm:"NOT NULL"`
	LoginType   LoginType
	LoginSource int64 `xorm:"NOT NULL DEFAULT 0"`
	LoginName   string
	Type        UserType
	Avatar      string

	IsAdmin bool

	OwnedOrgs   []*User       `xorm:"-" json:"-"`
	Orgs        []*User       `xorm:"-" json:"-"`
	Repos       []*Repository `xorm:"-" json:"-"`

	Roles []*Role                      `xorm:"-" json:"-"`
	DockerContainer []*DockerContainer `xorm:"-" json:"-"`
	DockerImage []*DockerImage         `xorm:"-" json:"-"`

	Location    string
	Website     string
	Rands       string `xorm:"VARCHAR(10)"`
	Salt        string `xorm:"VARCHAR(10)"`

	Created     time.Time `xorm:"-" json:"-"`
	CreatedUnix int64
	Updated     time.Time `xorm:"-" json:"-"`
	UpdatedUnix int64
}

func (user *User)GetByUser()(*User, error) {
	has,err := Engine.Get(user)
	if err != nil {
		return nil, err
	}
	if has {
		return user,nil
	} else {
		return nil, err
	}
}

func (user *User)Insert() error {
	_, err := Engine.Insert(user)
	if err != nil {
		return err
	}
	return nil
}

func (user *User)FindListByUser(page int, pageSize int, order string) (*[]User, error) {
	var allusers []User
	err := Engine.Limit(pageSize,page*pageSize).OrderBy(order).Find(&allusers) //Get id>3 limit 10 offset 20
	if err != nil {
		return nil, err
	}
	return &allusers, nil
}

func (user *User)Delete() error{
	_, err := Engine.Delete(user)
	return err
}

func (user *User) Update() error{
	affected, err := Engine.ID(user.Id).Update(user)

	if err != nil && affected != 0 {
		return err
	}

	return nil
}

func (usr *User)FindListByUserID() ([]*Repository, error) {
	var repoList [] *Repository
	err := Engine.Where("owner_id = ?", usr.Id).OrderBy("created_unix").Find(repoList)
	if err != nil {
		return nil, err
	}
	return repoList, nil
}