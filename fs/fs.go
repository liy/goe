package fs

import (
	"io"
)

type FileSystem interface {
	// pack and index
	Packs() ([]string, error)
	Pack(filename string) (io.ReadSeeker, error)
	PackIndex(filename string) (io.ReadCloser, error)

	// object
	Object(filename string) (io.ReadCloser, error)

	// reference
	Reference(referenceName string) (io.ReadCloser, error)
	ReferenceNames() []string
	PackedReference() (io.ReadCloser, error)
	
	// data
	Open(path string) (io.ReadCloser, error)
	ReadAll(io.ReadCloser) ([]byte, error)
}

