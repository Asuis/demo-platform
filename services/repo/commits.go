package repo

import (
	"container/list"
	"github.com/gogs/git-module"
)

func GetCommits(path string, page int, revision string) (*list.List, error){
	repo, err := git.OpenRepository("path")
	if err != nil {
		return nil,err
	}
	commits,err := repo.CommitsByRange(revision, page)
	if err != nil {
		return nil,err
	}
	return commits, nil
}

func MergeCommit() {}

func DelCommit() {}

func AddCommit() {}