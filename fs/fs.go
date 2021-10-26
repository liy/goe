package fs

import (
	"io"
	"os"
	"path/filepath"
)

type DotGit interface {
	PackIndex(filename string) (io.Reader, error)
	Pack(filename string) (io.ReadSeeker, error)
	Object(filename string) (io.Reader, error)
	File(path string) (io.Reader, error)
	Root() string
}

type FileSystem struct {
	root string
}

func NewFileSystem(path string) *FileSystem {
	root, _ := filepath.Abs(path)
	return &FileSystem{root}
}

func (fs FileSystem) Root() string {
	return fs.root
}

func (fs FileSystem) PackIndex(filename string) (io.Reader, error) {
	if filepath.IsAbs(filename) {
		return os.Open(filename)
	} 
	return os.Open(filepath.Join(fs.root, "objects", "pack" , filename))
}

func (fs FileSystem) Pack(filename string) (io.ReadSeeker, error) {
	if filepath.IsAbs(filename) {
		return os.Open(filename)
	} 
	return os.Open(filepath.Join(fs.root, "objects", "pack" , filename))
}

func (fs FileSystem) Object(hash string) (io.Reader, error) {
	return os.Open(filepath.Join(fs.root, "objects", hash[:2], hash[2:]))
}

func (fs FileSystem) File(path string) (io.Reader, error) {
	return os.Open(filepath.Join(fs.root, path))
}