package goe

import "github.com/liy/goe/object"

type Repository struct {
}

func (repo *Repository) Open(path string) error {
	return nil
}

// TODO: add sorting
func (repo *Repository) GetCommits() []object.Commit {
	return nil
}