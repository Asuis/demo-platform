package repo

import (
	"container/list"
	"github.com/gogs/git-module"
	"path"
)

func GetCommits(p string, page int, revision string) (*list.List, error){
	repo, err := git.OpenRepository(path.Join(BaseDir, p))
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