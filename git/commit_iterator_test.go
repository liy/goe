package git

import (
	"testing"

	"github.com/liy/goe/errors"
	"github.com/liy/goe/object"
	"github.com/liy/goe/plumbing"
	"github.com/liy/goe/tests"
)

func TestCommitIterator(t *testing.T) {
	repo, err := OpenRepository("../repos/topo-sort")
	if err != nil {
		t.Fatal(err)
	}

	refs := repo.GetReferences()

	tips := make([]*object.Commit, len(refs))
	for i, ref := range refs {
		var c *object.Commit

		raw, err := repo.ReadObject(repo.Peel(ref))
		if err != nil {
			t.Fatal(err)
		}

		// Get the real tagged object if the object is a annotated tag object
		if raw.Type == plumbing.OBJ_TAG {
			tag, err := object.DecodeTag(raw)
			if err != nil {
				t.Fatal(err)
			}

			raw, err = repo.ReadObject(tag.Target)
			if err != nil {
				t.Fatal(err)
			}
		}

		c, err = object.DecodeCommit(raw)
		if err != nil {
			// The raw object can be other object other than commit, simply ignore anything that is not a commit
			if err == errors.ErrRawObjectTypeWrong {
				continue
			}
		}

		tips[i] = c
	}

	var commits []*object.Commit
	itr := NewCommitIterator(repo, tips)
	for {
		c, err := itr.Next()
		if err != nil {
			if err == Done {
				break
			}
			t.Fatal(err)
		}
		commits = append(commits, c)
	}

	tests.ToMatchSnapshot(t, commits)
}