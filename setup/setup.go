package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/liy/goe/plumbing/indexfile"
	"github.com/liy/goe/tests"
	"github.com/liy/goe/utils"
)

func main() {
	var wg sync.WaitGroup
	fmt.Println("Downloading fixtures...")
	fs := tests.GetFixtures()
	for _, f := range fs {
		wg.Add(1)
		go func(f tests.RepositoryFixture) {
			defer wg.Done()
			download(f)
		}(f)
	}
	wg.Wait()
	
	SnapshotIndexFile()
}

func download(fixture tests.RepositoryFixture) {
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

func SnapshotIndexFile() {
	fixture := tests.GetFixture("topo-sort")
	file := fixture.IndexFile("./repos")
	idx, _ := indexfile.Decode(file)

	b, err := json.Marshal(idx)
	fmt.Println(err)
	err = os.WriteFile("./plumbing/indexfile/indexfile_snapshot.json", b, 0644)
	fmt.Println(err)
}
