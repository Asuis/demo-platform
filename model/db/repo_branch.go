package db

import "github.com/gogs/git-module"

type Branch struct {
	RepoPath string
	Name     string

	IsProtected bool
	Commit      *git.Commit
}
