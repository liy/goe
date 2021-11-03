package plumbing

import (
	"regexp"
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
	// local branch or remote branches which is not HEAD
	return strings.HasPrefix(refname, branchRef) || (IsRemote(refname) && !strings.HasSuffix(refname, "HEAD"))
}

func IsTag(refname string) bool {
	return strings.HasPrefix(refname, tagRef)
}

func IsHead(refname string) bool {
	mached, err := regexp.MatchString(`refs\/remotes\/\w+\/HEAD`, refname)
	if err != nil {
		return false
	}
	return refname == "HEAD" || mached
}

type ReferenceTarget string

func (rt ReferenceTarget) IsReference() bool {
	return strings.HasPrefix(string(rt), symbolicRef)
}

func (rt ReferenceTarget) IsHash() bool {
	return !rt.IsReference()
}

/*
ReferenceName can be a hash or another reference 
*/
func (rt ReferenceTarget) ReferenceName() string {
	if rt.IsReference() {
		return strings.TrimSpace(string(rt)[5:])
	}
	return strings.TrimSpace(string(rt))
}

func (rt ReferenceTarget) String() string {
	return string(rt)
}

var shorthandRegex = regexp.MustCompile(`refs/\w+/`)

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

func (r Reference) IsSymbolic() bool {
	return r.Target.IsReference()
}

func (r Reference) Shorthand() string {
	return shorthandRegex.ReplaceAllString(r.Name, "")
}