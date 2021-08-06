package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/liy/goe/src/goe/protobuf"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

func main() {
    start := time.Now()
	directory := "./repo"
	// Opens an already existing repository.
	r, err := git.PlainOpen(directory)
	CheckIfError(err)
	
	
	// ... retrieves the branch pointed by HEAD
	ref, err := r.Head()
	head := protobuf.Head {
		Hash: ref.Hash().String(),
		Name: ref.Name().String(),
		Shorthand: ref.Name().Short(),
	}
	CheckIfError(err)

	// ... retrieves the commit history
	since := time.Time{}
	until := time.Date(2021, 9, 30, 0, 0, 0, 0, time.UTC)
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash(), Since: &since, Until: &until, Order: git.LogOrderCommitterTime})
	CheckIfError(err)

    log.Printf("Log all commits took %s", time.Since(start))

	var commits []*protobuf.Commit
	err = cIter.ForEach(func(c *object.Commit) error {
		messages := strings.Split(c.Message, "\n")

		summary := messages[0]
		body := ""
		if len(messages) > 1 {
			body = strings.Join(messages[1:], "\n")
		}


		parents := make([]string, c.NumParents()) 
		for i, pc := range c.ParentHashes {
			parents[i] = pc.String()
		}

		commit := protobuf.Commit {
			Hash: c.Hash.String(),
			Summary: summary,
			Body: body,
			Author: &protobuf.Contact {
				Name: c.Author.Name,
				Email: c.Author.Email,
			},
			Committer:  &protobuf.Contact {
				Name: c.Committer.Name,
				Email: c.Committer.Email,
			},
			Parents: parents,
			CommitTime: ts.New(c.Committer.When),
		}
		commits = append(commits, &commit)

		// fmt.Println(commit)
		// fmt.Println("")

		return nil
	})
	CheckIfError(err)

	var references []*protobuf.Reference
	rIter, err := r.References()
	CheckIfError(err)

	err = rIter.ForEach(func(r *plumbing.Reference) error {
		ref := protobuf.Reference {
			Name: r.Name().String(),
			Shorthand: r.Name().Short(),
			Hash: r.Hash().String(),
			Type: protobuf.Reference_Type(r.Type()),
			IsRemote: r.Target().IsRemote(),
			IsBranch: r.Target().IsBranch(),
		}
		references = append(references, &ref)
		return nil
	})
	CheckIfError(err)


	data := protobuf.Repository {
		Id: directory,
		Commits: commits,
		References: references,
		Head: &head,
	}

	fmt.Println(data.Commits[0])

	CheckIfError(err)
}