package plumbing

import (
	"strings"
)

// branch, directly target to commit
// simple tag, direct target to object, recursive?
// annotated tag, target to tag object, then target to object recursive?
// HEAD can be referencing anything above, which means it points to another reference or directly points object
const (
	branchRef = "refs/heads/"
	tagRef    = "refs/tags/"
	noteRef   = "refs/notes/"
	remoteRef = "refs/remotes/"
)

const symbolicRef = "ref: "

func IsRemote(refname string) bool {
	return strings.HasPrefix(refname, remoteRef)
}

func IsBranch(refname string) bool {
	return strings.HasPrefix(refname, branchRef)
}

func IsTag(refname string) bool {
	return strings.HasPrefix(refname, tagRef)
}

func IsHead(refname string) bool {
	return refname == "HEAD"
}

type ReferenceTarget string

func (rt ReferenceTarget) IsRef() bool {
	return strings.HasPrefix(string(rt), symbolicRef)
}

func (rt ReferenceTarget) IsHash() bool {
	return !rt.IsRef()
}

/*
ReferenceName can be a hash or another reference 
*/
func (rt ReferenceTarget) ReferenceName() string {
	if rt.IsRef() {
		return strings.TrimSpace(string(rt)[5:])
	}
	return strings.TrimSpace(string(rt))
}

func (rt ReferenceTarget) String() string {
	return string(rt)
}

type Reference struct {
	Name   string
	Target ReferenceTarget
}

func NewReference(name string, target []byte) *Reference {
	return &Reference{
		Name:   name,
		Target: ReferenceTarget(target),
	}
}

// type ReferenceReader struct {
// 	path       string
// 	scanner    *bufio.Scanner
// 	refReaders map[string]*bufio.Reader
// }

// func NewReferenceReader(repoPath string) *ReferenceReader {
// 	var scanner *bufio.Scanner
// 	file, err := os.Open(filepath.Join(repoPath, "packed-refs"))
// 	if err != nil {
// 		scanner = bufio.NewScanner(file)
// 	}

// 	return &ReferenceReader{
// 		path:    repoPath,
// 		scanner: scanner,
// 	}
// }

// func (n *ReferenceReader) head() string {
// 	return filepath.Join(n.path, "HEAD")
// }

// func (n *ReferenceReader) Read(refname ReferenceName) []byte {
// 	if refname == "HEAD" {
// 		if reader, ok := r.refReaders[r.head()]; ok {
// 			return reader.ReadLine()
// 		}
// 	}
// }
