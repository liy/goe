package fs

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type DotGit struct {
	root string
}

func NewDotGit(root string) *DotGit {
	return &DotGit{root}
}

func (fs DotGit) Packs() ([]string, error) {
	return filepath.Glob(filepath.Join(fs.root, "objects/pack", "*.pack"))
}

func (fs DotGit) PackIndex(filename string) (io.ReadCloser, error) {
	return os.Open(path.Join(fs.root, "objects", "pack" , filename))
}

func (fs DotGit) Pack(filename string) (io.ReadSeeker, error) {
	return os.Open(path.Join(fs.root, "objects", "pack" , filename))
}

func (fs DotGit) Object(hash string) (io.ReadCloser, error) {
	return os.Open(path.Join(fs.root, "objects", hash[:2], hash[2:]))
}
func (fs DotGit) Reference(name string) (io.ReadCloser, error) {
	return fs.Open(name)
}

func (fs DotGit) PackedReference() (io.ReadCloser, error) {
	return fs.Open("packed-refs")
}

func (fs DotGit) Open(filepath string) (io.ReadCloser, error) {
	return os.Open(path.Join(fs.root, filepath))
}

func (fs DotGit) ReferenceNames() []string  {
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

func (fs DotGit) ReadAll(reader io.ReadCloser) ([]byte, error) {
	return ioutil.ReadAll(reader)
}