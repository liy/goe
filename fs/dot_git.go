package fs

import (
	"io"
)

type DotGit interface {
	Root() string
	Pack(filename string) (io.ReadSeeker, error)
	PackIndex(filename string) (io.ReadCloser, error)
	Object(filename string) (io.ReadCloser, error)
	Reference(referenceName string) (io.ReadCloser, error)
	ReferenceNames() []string
	PackedReference() (io.ReadCloser, error)
	Open(path string) (io.ReadCloser, error)
}
