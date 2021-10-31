package main

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/liy/goe/utils"
)

type Download struct {
	Name         string
	Url          string
	PackfileHash string
}

var downloads = []Download{
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

func download(fixture Download) {
	defer fmt.Printf("%v downloaded\n", fixture.Name)

	var containerFolder = "./repos/"
	p := filepath.Join(containerFolder, fixture.Name)
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		return
	}

	// Get the data
	res, err := http.Get(fixture.Url)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	gzipReader, err := gzip.NewReader(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	defer gzipReader.Close()

	err = utils.Untar(gzipReader, containerFolder)
	if err != nil {
		fmt.Println(err)
	}
}

func GetDownloads() []Download {
	return downloads
}

func GetDownload(name string) *Download {
	for _, f := range downloads {
		if f.Name == name {
			return &f
		}
	}

	return nil
}

func main() {
	var wg sync.WaitGroup
	fmt.Println("Downloading fixtures...")
	fs := GetDownloads()
	for _, f := range fs {
		wg.Add(1)
		go func(f Download) {
			defer wg.Done()
			download(f)
		}(f)
	}
	wg.Wait()
}
