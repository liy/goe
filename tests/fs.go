package tests

import (
	"io"
	"io/ioutil"
	"path"
	"strings"
)

type EmbededDotGit struct {
}

func NewEmbededDotGit() *EmbededDotGit {
	return &EmbededDotGit{}
}

func (fs EmbededDotGit) Packs() ([]string, error) {
	file, _ := FS(false).Open("/topo-sort/.git/objects/pack")
	list, _ := file.Readdir(-1)
	packNames := make([]string, 0)
	for _, f := range list {
		if strings.HasSuffix(f.Name(), ".idx") {
			continue
		}
		packNames  = append(packNames, f.Name())
	}

	return packNames, nil
}

func (fs EmbededDotGit) PackIndex(filename string) (io.ReadCloser, error) {
	return FS(false).Open(path.Join("/topo-sort/.git", "objects", "pack", filename))
}

func (fs EmbededDotGit) Pack(filename string) (io.ReadSeeker, error) {
	return FS(false).Open(path.Join("/topo-sort/.git", "objects", "pack", filename))
}

func (fs EmbededDotGit) Object(hash string) (io.ReadCloser, error) {
	return FS(false).Open(path.Join("/topo-sort/.git", "objects", hash[:2], hash[2:]))
}

func (fs EmbededDotGit) Reference(name string) (io.ReadCloser, error) {
	return FS(false).Open(path.Join("/topo-sort/.git", name))
}

func (fs EmbededDotGit) PackedReference() (io.ReadCloser, error) {
	return FS(false).Open(path.Join("/topo-sort/.git", "packed-refs"))
}

func (fs EmbededDotGit) Open(filepath string) (io.ReadCloser, error) {
	return FS(false).Open(path.Join(filepath))
}

func (fs EmbededDotGit) ReferenceNames() []string {
	var names []string

	var traverse func(string)
	traverse = func(p string) {
		file, _ := FS(false).Open(path.Join("/topo-sort/.git/", p))
		entries, err := file.Readdir(-1)
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

func (fs EmbededDotGit) ReadAll(reader io.ReadCloser) ([]byte, error) {
	return ioutil.ReadAll(reader)
}