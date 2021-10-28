package tests

import (
	"fmt"
	"os"
	"path/filepath"
)

type RepositoryFixture struct {
	Name         string
	Url          string
	PackfileHash string
}

var fixtures = []RepositoryFixture{
	{"nodegit", "https://gitlab.com/liyss/goe-fixtures/-/raw/main/nodegit.tar.gz", "46826dbf994f5bebc4a8ac318ab631b17357659d"},
	{"rails", "https://gitlab.com/liyss/goe-fixtures/-/raw/main/rails.tar.gz", "318d97c6cd13eaf0ce20677d1b667073092b1bb1"},
	{
		"git",
		"https://gitlab.com/liyss/goe-fixtures/-/raw/main/git.tar.gz",
		"0829575a095e1d5bb1c752dd60d45568a049ed76",
	},
	{
		"topo-sort",
		"https://gitlab.com/liyss/goe-fixtures/-/raw/main/topo-sort.tar.gz",
		"6179faab20f2d649a12fd52aab3c8d6e32b27dcd",
	},
}


func GetFixtures() []RepositoryFixture {
	return fixtures
}

func GetFixture(name string) *RepositoryFixture {
	for _, f := range fixtures {
		if f.Name == name {
			return &f
		}
	}

	return nil
}

func (f *RepositoryFixture) GetIndexFilePath(folder string) string {
	p := fmt.Sprintf("./%s/.git/objects/pack/pack-%s.idx", f.Name, f.PackfileHash)
	return filepath.Join(folder, p)
}

func (f *RepositoryFixture) IndexFile(folder string) *os.File {
	file, err := os.Open(f.GetIndexFilePath(folder))
	if err != nil {
		fmt.Println(err)
	}
	return file
}

func (f *RepositoryFixture) GetPackFilePath(folder string) string {
	p := fmt.Sprintf("./%s/.git/objects/pack/pack-%s.pack", f.Name, f.PackfileHash)
	return filepath.Join(folder, p)
}

func (f *RepositoryFixture) PackFile(folder string) *os.File {
	file, err := os.Open(f.GetPackFilePath(folder))
	if err != nil {
		fmt.Println(err)
	}
	return file
}