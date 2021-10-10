package plumbing

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