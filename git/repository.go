package goe

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	goeErrors "github.com/liy/goe/errors"
	"github.com/liy/goe/object"
	"github.com/liy/goe/plumbing"
	"github.com/liy/goe/plumbing/packfile"
	"github.com/liy/goe/utils"
)

type Repository struct {
	path string
	packReaders []packfile.PackReader
}

func OpenRepository(path string) (repo *Repository, err error) {
	if path, err = filepath.Abs(filepath.Join(path, ".git")); err != nil {
		return nil, err
	}

	packFiles, err := filepath.Glob(filepath.Join(path, "objects/pack", "*.pack"))
	if err != nil {
		return nil, err
	}

	var readers []packfile.PackReader = make([]packfile.PackReader, len(packFiles))
	for i, p := range packFiles {
		readers[i] = *packfile.NewPackReader(p)
	}

	return &Repository {
		path,
		readers,
	}, nil
}

func (r *Repository) ReadObject(hash plumbing.Hash) (raw *plumbing.RawObject, err error) {
	for _, pr := range r.packReaders {
		raw, err = pr.ReadObject(hash)
		if err != nil {
			return nil, err
		}
	}
	return raw, nil
}

func (r *Repository) GetCommit(hash plumbing.Hash) (c *object.Commit, err error) {
	for _, pr := range r.packReaders {
		raw, err := pr.ReadObject(hash)
		if err == goeErrors.ErrObjectNotFound {
			continue
		} else if err != nil {
			return nil, err
		}

		return object.DecodeCommit(raw)
	}

	obj, err := object.ParseObjectFile(hash, r.path)
	if err != nil && err != goeErrors.ErrObjectNotFound {
		return nil, err
	}

	if obj.Type != "commit" {
		return nil, fmt.Errorf("object is not a commit")
	}

	c = &object.Commit{
		Hash: hash,
	}
	err = c.Decode(obj.Data)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (r *Repository) GetCommits(hash plumbing.Hash) ([]*object.Commit, error) {
	start, err := r.GetCommit(hash)
	if err != nil {
		return nil, err
	}
	queue := utils.NewPrioQueue(start)

	// count := 0
	commits := make([]*object.Commit, 0)
	visited :=  make(map[plumbing.Hash]bool)

	for {
		commit, ok := (*queue.Dequeue()).(*object.Commit)
		if !ok {
			return nil, errors.New("not a commit object")
		}
		commits = append(commits, commit)

		for _, ph := range commit.Parents {
			if _, exist := visited[ph]; !exist {
				visited[ph] = true
				p, err := r.GetCommit(ph)
				if err != nil {
					return nil, errors.New("cannot get parent commit: " + ph.String()) 
				}

				queue.Enqueue(p)
			}
		}

		if queue.Size() == 0 {
			break
		}
	}

	return commits, nil
}

func (r *Repository) GetTag(hash plumbing.Hash) (*object.Tag, error) {
	for _, pr := range r.packReaders {
		raw, err := pr.ReadObject(hash)
		if err == goeErrors.ErrObjectNotFound {
			continue
		} else if err != nil {
			return nil, err
		}

		return object.DecodeTag(raw)
	}

	obj, err := object.ParseObjectFile(hash, r.path)
	if err != nil {
		return nil, err
	}

	if obj.Type != "tag" {
		return nil, fmt.Errorf("object is not a tag")
	}

	t := &object.Tag{
		Hash: hash,
	}
	err = t.Decode(obj.Data)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *Repository) GetPath() string {
	return r.path
}

// TODO: cache references?
func (r *Repository) findPackedReference(scanner *bufio.Scanner, name string) (plumbing.Hash, error) {
	for scanner.Scan() {
		lineBytes := scanner.Bytes()
		chunks := bytes.SplitN(lineBytes, []byte{' '}, 2)
		
		if name == string(chunks[1]) {
			return plumbing.NewHash(chunks[0]), nil
		}
	}

	return plumbing.Hash{}, goeErrors.ErrReferenceNoteFound
}

func (r *Repository) GetReference(name string) (plumbing.Hash, error) {
	// Try packed-refs first
	file, err := os.Open(filepath.Join(r.path, "packed-refs"))
	if err == nil {
		return r.findPackedReference(bufio.NewScanner(file), name)
	}

	// then refs folder
	lineBytes, err := ioutil.ReadFile(filepath.Join(r.path, name))
	if err != nil {
		return plumbing.Hash{}, goeErrors.ErrReferenceNoteFound
	}
	return plumbing.NewHash(bytes.TrimSpace(lineBytes)), nil
}

func (r *Repository) Head() (plumbing.Hash, error) {
	file, err := os.Open(filepath.Join(r.path, "HEAD"))
	if err != nil {
		return plumbing.Hash{}, goeErrors.ErrReferenceNoteFound
	}
	reader := bufio.NewReader(file)
	line, _, _ := reader.ReadLine()
	refNameBytes := bytes.SplitN(line, []byte{' '}, 2)[1]
	refNameBytes = bytes.TrimSpace(refNameBytes)
	return r.GetReference(string(refNameBytes))
}

func (r *Repository) GetReferences() ([]plumbing.Reference, error) {
	references := make([]plumbing.Reference, 0)
	// Try packed-refs first
	file, err := os.Open(filepath.Join(r.path, "packed-refs"))
	if err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lineBytes := scanner.Bytes()
			
			// Ignore comments and annotated tag
			if lineBytes[0] == '#' || lineBytes[0] == '^' {
				continue
			}
	
			chunks := bytes.SplitN(lineBytes, []byte{' '}, 2)
			references = append(references, *plumbing.NewReference(string(chunks[1]), plumbing.ToHash(string(chunks[0]))))
		}
		return references, nil
	}
	
	// read individual ref files
	folder, err := os.Open(filepath.Join(r.path + "heads"))
    if err != nil {
        return nil, err
    }
    infos, err := folder.Readdir(0)
    if err != nil {
        return nil, err
    }

    for _, f := range infos {
		if f.IsDir() {
			continue
		}

		refName := "refs/heads/" + f.Name()
		lineBytes, err := os.ReadFile(filepath.Join(r.path, refName))
		if err != nil {
			continue
		}

		references = append(references, *plumbing.NewReference(refName, plumbing.ToHash(string(lineBytes))))
	}
	return references, nil
}