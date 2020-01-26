package repo

import (
	"github.com/gogs/git-module"
)

func SearchDir(path string,) (git.Entries, error){
	repo, err := git.OpenRepository(path)
	if err != nil {
		return nil, err
	}
	branch, err:=repo.GetHEADBranch()
	if err != nil {
		return nil, err
	}
	commit,err:=repo.GetBranchCommit(branch.Name)
	if err != nil {
		return nil, err
	}
	tree, err := repo.GetTree(commit.ID.String())
	if err != nil {
		return nil, err
	}
	entries, err := tree.ListEntries()
	if err != nil {
		return nil, err
	}
	return entries,nil
}

func GetRawFile(path string, relpath string) (*git.Blob, error){
	repo, err := git.OpenRepository(path)
	if err != nil {
		return nil, err
	}
	branch, err:=repo.GetHEADBranch()
	if err != nil {
		return nil, err
	}
	commit,err:=repo.GetBranchCommit(branch.Name)
	if err != nil {
		return nil, err
	}
	tree, err := repo.GetTree(commit.ID.String())
	if err != nil {
		return nil, err
	}

	blob, err := tree.GetBlobByPath(relpath)
	return blob,err
}

