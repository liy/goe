package goe

import (
	"fmt"
	"path/filepath"

	"github.com/liy/goe/errors"
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

func (r *Repository) ReadObject(hash plumbing.Hash) (raw *plumbing.RawObject, err error) {
	for _, pr := range r.packReaders {
		raw, err = pr.ReadObject(hash)
		if err != nil {
			return nil, err
		}
	}
	return raw, nil
}

// TODO: add sorting
func (repo *Repository) GetCommits() []object.Commit {
	return nil
}

func (r *Repository) GetCommit(hash plumbing.Hash) (c *object.Commit, err error) {
	for _, pr := range r.packReaders {
		raw, err := pr.ReadObject(hash)
		if err == errors.ErrObjectNotFound {
			continue
		} else if err != nil {
			return nil, err
		}

		return object.DecodeCommit(raw)
	}

	obj, err := object.ParseObjectFile(hash, r.path)
	if err != nil && err != errors.ErrObjectNotFound {
		return nil, err
	}

	if obj.Type != "commit" {
		return nil, fmt.Errorf("object is not a commit")
	}

	c = &object.Commit{
		Hash: hash,
	}
	err = c.Decode(obj.Data)
	if err != nil {
		return nil, err
	}

	return c, nil
}


func (r *Repository) GetTag(hash plumbing.Hash) (*object.Tag, error) {
	for _, pr := range r.packReaders {
		raw, err := pr.ReadObject(hash)
		if err == errors.ErrObjectNotFound {
			continue
		} else if err != nil {
			return nil, err
		}

		return object.DecodeTag(raw)
	}

	obj, err := object.ParseObjectFile(hash, r.path)
	if err != nil {
		return nil, err
	}

	if obj.Type != "tag" {
		return nil, fmt.Errorf("object is not a tag")
	}

	t := &object.Tag{
		Hash: hash,
	}
	err = t.Decode(obj.Data)
	if err != nil {
		return nil, err
	}

	return t, nil
}

