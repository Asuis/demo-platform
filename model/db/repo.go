package db

import (
	"time"
)

type Repository struct {
	ID int64
	OwnerID         int64  `xorm:"UNIQUE(s)"`
	Owner           *User  `xorm:"-" json:"-"`
	LowerName       string `xorm:"UNIQUE(s) INDEX NOT NULL"`
	Name            string `xorm:"INDEX NOT NULL"`
	Description     string `xorm:"VARCHAR(512)"`
	Website         string
	DefaultBranch   string
	Size            int64 `xorm:"NOT NULL DEFAULT 0"`
	UseCustomAvatar bool

	// Counters
	NumWatches          int
	NumStars            int
	NumForks            int
	NumIssues           int
	NumClosedIssues     int
	NumOpenIssues       int `xorm:"-" json:"-"`
	NumPulls            int
	NumClosedPulls      int
	NumOpenPulls        int `xorm:"-" json:"-"`
	NumMilestones       int `xorm:"NOT NULL DEFAULT 0"`
	NumClosedMilestones int `xorm:"NOT NULL DEFAULT 0"`
	NumOpenMilestones   int `xorm:"-" json:"-"`
	NumTags             int `xorm:"-" json:"-"`

	IsPrivate bool
	IsBare    bool

	IsMirror bool
	*Mirror  `xorm:"-" json:"-"`

	// Advanced settings
	EnableWiki            bool `xorm:"NOT NULL DEFAULT true"`
	AllowPublicWiki       bool
	EnableExternalWiki    bool
	ExternalWikiURL       string
	EnableIssues          bool `xorm:"NOT NULL DEFAULT true"`
	AllowPublicIssues     bool
	EnableExternalTracker bool
	ExternalTrackerURL    string
	ExternalTrackerFormat string
	ExternalTrackerStyle  string
	ExternalMetas         map[string]string `xorm:"-" json:"-"`
	EnablePulls           bool              `xorm:"NOT NULL DEFAULT true"`
	PullsIgnoreWhitespace bool              `xorm:"NOT NULL DEFAULT false"`
	PullsAllowRebase      bool              `xorm:"NOT NULL DEFAULT false"`

	IsFork   bool `xorm:"NOT NULL DEFAULT false"`
	ForkID   int64
	BaseRepo *Repository `xorm:"-" json:"-"`

	Created     time.Time `xorm:"-" json:"-"`
	CreatedUnix int64
	Updated     time.Time `xorm:"-" json:"-"`
	UpdatedUnix int64
}

func CreateRepository(repository * Repository) error {
	_, err := Engine.Insert(repository)
	if err != nil {
		return err
	}
	return nil
}

func GetRepository(repository *Repository) (bool, error) {
	has, err := Engine.Get(repository)
	if err !=nil {
		return false, nil
	}

	return has, nil
}

func ListRepository(repository *Repository, page int, pageSize int, order string) (*[]Repository, error){

	var res []Repository

	err := Engine.Where("OwnerId = ?", repository.OwnerID).Limit(pageSize,page*pageSize).OrderBy(order).Find(&res) //Get id>3 limit 10 offset 20

	if err != nil {
		return nil, err
	}

	return &res, nil
}
