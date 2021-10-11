package plumbing

const (
	branchPath = "refs/heads/"
	tagPath    = "refs/tags/"
	remotePath = "refs/remotes/"
)

type Reference struct {
	Name string
	Hash Hash
}

func NewReference(name string, hash Hash) *Reference {
	return &Reference{
		name,
		hash,
	}
}

func ReadBranches(name string, repoPath string) {

}

func ReadTags(name string, repoPath string) {

}

func ReadRemoteBranches(name string, repoPath string) {

}

func ReadRemoteTags(name string, repoPath string) {

}