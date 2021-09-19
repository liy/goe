package goe

import (
	"fmt"
	"path/filepath"

	"github.com/liy/goe/object"
	"github.com/liy/goe/plumbing"
	"github.com/liy/goe/plumbing/packfile"
)

type Repository struct {
	path string
	packReaders []packfile.PackReader
}

func OpenRepository(path string) (repo *Repository, err error) {
	if path, err = filepath.Abs(filepath.Join(path, ".git")); err != nil {
		return nil, err
	}

	packFiles, err := filepath.Glob(filepath.Join(path, "objects/pack", "*.pack"))
	if err != nil {
		return nil, err
	}

	var readers []packfile.PackReader = make([]packfile.PackReader, len(packFiles))
	for i, p := range packFiles {
		readers[i] = *packfile.NewPackReader(p)
	}

	return &Repository {
		path,
		readers,
	}, nil
}

// TODO: add sorting
func (repo *Repository) GetCommits() []object.Commit {
	return nil
}

func (r *Repository) GetCommit(hash plumbing.Hash) (c *object.Commit, err error) {
	for _, pr := range r.packReaders {
		raw, err := pr.ReadObject(hash)
		if err != nil {
			// TODO: find it in the file system
			return nil, err
		} else {
			return object.DecodeCommit(raw)
		}
	}

	return nil, fmt.Errorf("cannot find commit %s", hash)
}