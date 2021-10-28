package git

import (
	"bufio"
	"bytes"
	"path/filepath"
	"regexp"

	"github.com/liy/goe/errors"
	"github.com/liy/goe/fs"
	"github.com/liy/goe/object"
	"github.com/liy/goe/plumbing"
	"github.com/liy/goe/plumbing/packfile"
)

var isDotGit, _ = regexp.Compile(".git/?$")

type Repository struct {
	path string
	dotgit		fs.FileSystem
	packReaders []packfile.PackReader
}

func OpenRepository(path string) (repo *Repository, err error) {
	if !isDotGit.MatchString(path) {
		path = filepath.Join(path, ".git")
	}
	dotgit := fs.NewDotGit(path)

	packs, err := dotgit.Packs()
	if err != nil {
		return nil, err
	}

	var readers []packfile.PackReader = make([]packfile.PackReader, len(packs))
	for i, p := range packs {
		packName := filepath.Base(p)
		pack, err := dotgit.Pack(packName)
		if err != nil {
			return nil, err
		}
		idx, err := dotgit.PackIndex(packName[:len(packName)-4] + "idx")
		if err != nil {
			return nil, err
		}

		readers[i] = *packfile.NewPackReader(pack, idx)
	}

	return &Repository{
		path,
		dotgit,
		readers,
	}, nil
}

func (r *Repository) ReadObjectFile(hash plumbing.Hash) (*plumbing.RawObject, error) {
	file, err := r.dotgit.Object(hash.String())
	if err != nil {
		return nil, err
	}

	raw := plumbing.NewRawObject(hash)
	err = raw.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

func (r *Repository) ReadObject(hash plumbing.Hash) (raw *plumbing.RawObject, err error) {
	// Try find it in all packed files
	for _, pr := range r.packReaders {
		raw, err = pr.ReadObject(hash)
		if err == errors.ErrObjectNotFound {
			continue
		} else if err != nil {
			return nil, err
		}

		return raw, nil
	}

	// Try find it in object folder
	raw, err = r.ReadObjectFile(hash)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

func (r *Repository) GetCommit(hash plumbing.Hash) (c *object.Commit, err error) {
	var raw *plumbing.RawObject
	for _, pr := range r.packReaders {
		raw, err = pr.ReadObject(hash)
		if err == errors.ErrObjectNotFound {
			continue
		} else if err != nil {
			return nil, err
		}

		return object.DecodeCommit(raw)
	}

	// Try find it in object folder
	raw, err = r.ReadObjectFile(hash)
	if err != nil {
		return nil, err
	}

	return object.DecodeCommit(raw)
}

func (r *Repository) GetCommits(hash plumbing.Hash) ([]*object.Commit, error) {
	start, err := r.GetCommit(hash)
	if err != nil {
		return nil, err
	}
	queue := CommitPrioQueue{}
	queue.Enqueue(start)

	// count := 0
	commits := make([]*object.Commit, 0)
	visited := make(map[plumbing.Hash]bool)

	for {
		commit := queue.Dequeue()
		commits = append(commits, commit)

		for _, ph := range commit.Parents {
			if _, exist := visited[ph]; !exist {
				visited[ph] = true
				p, err := r.GetCommit(ph)
				if err != nil {
					return nil, object.NewParentError(ph)
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

func (r *Repository) GetAnnotatedTag(hash plumbing.Hash) (t *object.Tag, err error) {
	var raw *plumbing.RawObject
	for _, pr := range r.packReaders {
		raw, err = pr.ReadObject(hash)
		if err == errors.ErrObjectNotFound {
			continue
		} else if err != nil {
			return nil, err
		}

		return object.DecodeTag(raw)
	}

	// Try find it in object folder
	raw, err = r.ReadObjectFile(hash)
	if err != nil {
		return nil, err
	}

	return object.DecodeTag(raw)
}

func (r *Repository) GetPath() string {
	return r.path
}

func (r *Repository) GetReferences() []*plumbing.Reference {
	refs := make([]*plumbing.Reference, 0)

	// record all packed references
	packed := make(map[string]bool)
	file, err := r.dotgit.PackedReference()
	if err == nil {
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lineBytes := scanner.Bytes()

			// Ignore comments and annotated tag
			if lineBytes[0] == '#' || lineBytes[0] == '^' {
				continue
			}

			chunks := bytes.SplitN(lineBytes, []byte{' '}, 2)
			refs = append(refs, plumbing.NewReference(string(chunks[1]), chunks[0]))
			packed[string(chunks[1])] = true
		}
	}

	// Find all loose references
	names := r.dotgit.ReferenceNames()
	for _, n := range names {
		if packed[n] {
			continue
		}
		file, _ := r.dotgit.Reference(n)
		data, _ := r.dotgit.ReadAll(file)
		refs = append(refs, plumbing.NewReference(n, bytes.TrimSpace(data)))
	}

	// Head
	ref, err := r.HEAD()
	if err == nil {
		refs = append(refs, ref)
	}

	return refs
}

// symbolic reference which reference to another reference
func (r *Repository) HEAD() (*plumbing.Reference, error) {
	file, err := r.dotgit.Open("HEAD")
	if err != nil {
		return nil, errors.ErrReferenceNotFound
	}
	
	data, err := r.dotgit.ReadAll(file)
	if err != nil {
		return nil, errors.ErrReferenceNotFound
	}

	return plumbing.NewReference("HEAD", bytes.TrimSpace(data)), nil
}

func (r *Repository) Peel(reference *plumbing.Reference) plumbing.Hash {
	if reference.Target.IsHash() {
		return plumbing.ToHash(reference.Target.String())
	}

	// symbolic reference points to another ref
	ref, err := r.GetReference(reference.Target.ReferenceName())
	// the pointed reference does not exist
	if err != nil {
		return plumbing.ZeroHash
	}

	if ref.Target.IsReference() {
		return r.Peel(ref)
	}

	return plumbing.ToHash(ref.Target.String())
}

func (r *Repository) GetBranch(name string) (*plumbing.Reference, error) {
	return r.GetReference("refs/heads/" + name)
}

func (r *Repository) GetRemoteBranch(branchName string, remoteName string) (*plumbing.Reference, error) {
	return r.GetReference("refs/remotes/" + remoteName + "/heads" + branchName)
}

func (r *Repository) GetReference(refname string) (*plumbing.Reference, error) {
	// Try packed-refs first
	file, err := r.dotgit.PackedReference()
	if err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lineBytes := scanner.Bytes()

			// Ignore comments and annotated tag
			if lineBytes[0] == '#' || lineBytes[0] == '^' {
				continue
			}

			chunks := bytes.SplitN(lineBytes, []byte{' '}, 2)

			if refname == string(chunks[1]) {
				return plumbing.NewReference(refname, chunks[0]), nil
			}
		}
	}

	// then refs folder
	file, err = r.dotgit.Reference(refname)
	if err != nil {
		return nil, errors.ErrReferenceNotFound
	}
	lineBytes, err := r.dotgit.ReadAll(file)
	if err != nil {
		return nil, errors.ErrReferenceNotFound
	}
	return plumbing.NewReference(refname, lineBytes), nil
}