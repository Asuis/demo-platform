package repo

import (
	"container/list"
	"github.com/gogs/git-module"
	"path"
	"strings"
)

func SearchDir(p string, relpath string) (*git.Entries, error){
	if !strings.HasSuffix(p, ".git") {
		p += ".git"
	}
	repo, err := git.OpenRepository(path.Join(BaseDir, p))
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
	tr, err := tree.SubTree(relpath)
	if err != nil {
		return nil, err
	}
	entries, err := tr.ListEntries()
	if err != nil {
		return nil, err
	}

	return &entries,nil
}

func GetRawFile(p string, relpath string) (*git.Blob, error){
	if !strings.HasSuffix(p, ".git") {
		p += ".git"
	}
	repo, err := git.OpenRepository(path.Join(BaseDir, p))
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

func GetCommit(path string, relpath string, page int, pageSize int)(*list.List, error) {
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
	li, err := commit.CommitsByRangeSize(page, pageSize)
	if err != nil {
		return nil, err
	}
	return li, nil
}