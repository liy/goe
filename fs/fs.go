package fs

import (
	"io"
	"os"
	"path"
	"regexp"
)

type FileSystem struct {
	root string
}


var isDotGit, _ = regexp.Compile(".git/?$")

func NewFileSystem(root string) *FileSystem {
	if isDotGit.MatchString(root) {
		return &FileSystem{root}
	}
	return &FileSystem{path.Join(root, ".git")}
}

func (fs FileSystem) Root() string {
	return fs.root
}

func (fs FileSystem) PackIndex(filename string) (io.ReadCloser, error) {
	return os.Open(path.Join(fs.root, "objects", "pack" , filename))
}

func (fs FileSystem) Pack(filename string) (io.ReadSeeker, error) {
	return os.Open(path.Join(fs.root, "objects", "pack" , filename))
}

func (fs FileSystem) Object(hash string) (io.ReadCloser, error) {
	return os.Open(path.Join(fs.root, "objects", hash[:2], hash[2:]))
}
func (fs FileSystem) Reference(name string) (io.ReadCloser, error) {
	return fs.Open(name)
}

func (fs FileSystem) PackedReference() (io.ReadCloser, error) {
	return fs.Open("packed-refs")
}

func (fs FileSystem) Open(filepath string) (io.ReadCloser, error) {
	return os.Open(path.Join(fs.root, filepath))
}

func (fs FileSystem) ReferenceNames() []string  {
	var names []string

	var traverse func(string)
	traverse = func(p string) {
		entries, err := os.ReadDir(path.Join(fs.root, p))
		if err != nil {
			panic(err)
		}
		
		for _, entry := range entries {
			if entry.IsDir() {
				traverse(path.Join(p, entry.Name()))
			} else {
				names = append(names, path.Join(p, entry.Name()))
			}
		}
	}

	traverse("refs")

	return names
}