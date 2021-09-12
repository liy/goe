package object

import "github.com/liy/goe/plumbing"

type Commit struct {
	Hash plumbing.Hash
}

func ReadRaw(raw *plumbing.RawObject) error {
	return nil
}