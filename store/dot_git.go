package store

import (
	"net/http"
	"path"
	"strings"
)

type DotGit struct {
	Root string
	fs http.FileSystem
}

func NewDotGit(root string, fs http.FileSystem) *DotGit {
	return &DotGit{root, fs}
}

func (d DotGit) Open(filepath string) (http.File, error) {
	return d.fs.Open(path.Join(d.Root, filepath))
}

func (d DotGit) Pack(hash string) (http.File, error) {
	return d.Open(path.Join("objects", "pack" , "pack-"+ hash + ".pack"))
}

func (d DotGit) PackIndex(hash string) (http.File, error) {
	return d.Open(path.Join("objects", "pack" , "pack-"+ hash + ".idx"))
}

func (d DotGit) Packs(callback func(pReader http.File, idxReader http.File)) {
	folder, _ := d.Open("objects/pack")
	entries, _ := folder.Readdir(-1)
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".pack") {
			packName := entry.Name()
			indexName := packName[:len(packName)-4] + "idx"
			pack, err := d.Open(path.Join("objects/pack", packName))
			if err != nil {
				panic(err)
			}
			idx, err := d.Open(path.Join("objects/pack", indexName))
			if err != nil {
				panic(err)
			}

			callback(pack, idx)
		}
	}
}

func (d DotGit) Object(hash string) (http.File, error) {
	return d.Open(path.Join("objects", hash[:2], hash[2:]))
}

func (d DotGit) Reference(name string) (http.File, error) {
	return d.Open(path.Join(name))
}

func (d DotGit) PackedReference() (http.File, error) {
	return d.Open(path.Join("packed-refs"))
}

func (d DotGit) ReferenceNames() []string  {
	var names []string

	var traverse func(string)
	traverse = func(p string) {
		file, err := d.Open(p)
		if err != nil {
			panic(err)
		}
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