package object

import "github.com/liy/goe/plumbing"

type Object interface {
	Parse(raw *plumbing.RawObject) error
}