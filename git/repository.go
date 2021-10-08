package goe

import (
	"errors"
	"fmt"
	"path/filepath"

	goeErrors "github.com/liy/goe/errors"
	"github.com/liy/goe/object"
	"github.com/liy/goe/plumbing"
	"github.com/liy/goe/plumbing/packfile"
	"github.com/liy/goe/utils"
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

func (r *Repository) GetCommit(hash plumbing.Hash) (c *object.Commit, err error) {
	for _, pr := range r.packReaders {
		raw, err := pr.ReadObject(hash)
		if err == goeErrors.ErrObjectNotFound {
			continue
		} else if err != nil {
			return nil, err
		}

		return object.DecodeCommit(raw)
	}

	obj, err := object.ParseObjectFile(hash, r.path)
	if err != nil && err != goeErrors.ErrObjectNotFound {
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

func (r *Repository) GetCommits(hash plumbing.Hash) ([]*object.Commit, error) {
	start, err := r.GetCommit(hash)
	if err != nil {
		return nil, err
	}
	queue := utils.NewPrioQueue(start)

	// count := 0
	commits := make([]*object.Commit, 0)
	visited :=  make(map[plumbing.Hash]bool)

	for {
		commit, ok := (*queue.Dequeue()).(*object.Commit)
		if !ok {
			return nil, errors.New("not a commit object")
		}
		commits = append(commits, commit)

		for _, ph := range commit.Parents {
			if _, exist := visited[ph]; !exist {
				visited[ph] = true
				p, err := r.GetCommit(ph)
				if err != nil {
					return nil, errors.New("cannot get parent commit: " + ph.String()) 
				}

				queue.Enqueue(p)
			}
		}

		if queue.Size() == 0 {
			break
		}
	}

	return commits, nil
}

func (r *Repository) GetTag(hash plumbing.Hash) (*object.Tag, error) {
	for _, pr := range r.packReaders {
		raw, err := pr.ReadObject(hash)
		if err == goeErrors.ErrObjectNotFound {
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

