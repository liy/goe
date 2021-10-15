package plumbing

import "strings"

const (
	branchPath = "refs/heads/"
	tagPath    = "refs/tags/"
	remotePath = "refs/remotes/"
	notePath   = "refs/notes/"
)

type Reference struct {
	Name string
	target Hash
}

func NewReference(name string, target Hash) *Reference {
	return &Reference{
		name,
		target,
	}
}

func (r *Reference) IsRemote() bool {
	return strings.HasPrefix(r.Name, remotePath)
}

func (r *Reference) IsBranch() bool {
	return strings.HasPrefix(r.Name, branchPath)
}

func (r *Reference) IsAnnotatedTag() bool {

}

func (r *Reference) Target() Hash {
	return r.target
}