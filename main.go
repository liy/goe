package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/liy/goe/git"
	"github.com/liy/goe/object"
	"github.com/liy/goe/plumbing"

	"github.com/liy/goe/src/protobuf"
	ts "google.golang.org/protobuf/types/known/timestamppb"

	_ "net/http/pprof"

	_ "github.com/pkg/profile"
)


func mine() {
	start := time.Now()

	r, err := git.SimpleOpen("./repos/rails")
	if err != nil {
		fmt.Println(err)
	}

	refs := r.GetReferences()

	// Setup potential tips
	tips := make([]*object.Commit, len(refs))
	for i, ref := range refs {
		var c *object.Commit

		raw, err := r.ReadObject(r.Peel(ref))
		if err != nil {
			fmt.Println(err)
			continue
		}

		if raw.Type == plumbing.OBJ_TAG {
			tag, err := object.DecodeTag(raw)
			if err != nil {
				fmt.Println(err)
				continue
			}

			raw, err = r.ReadObject(tag.Target)
			if err != nil {
				continue
			}
		}

		c, err = object.DecodeCommit(raw)
		if err != nil {
			fmt.Println(err)
			continue
		}

		tips[i] = c
	}

	var commits []*protobuf.Commit
	cItr := git.NewCommitIterator(r, tips)
	for {
		c, err := cItr.Next()
		if err == git.Done {
			break
		}
		if err != nil {
			break
		}

		parents := make([]string, len(c.Parents))
		for i, ph := range c.Parents {
			parents[i] = ph.String()
		}

		chunks := strings.Split(c.Message, "\n")
		body := ""
		if len(chunks) == 2 {
			body = chunks[1]
		}
		commits = append(commits, &protobuf.Commit{
			Hash:    c.Hash.String(),
			Summary: chunks[0],
			Body:    body,
			Author: &protobuf.Contact{
				Name:  c.Author.Name,
				Email: c.Author.Email,
			},
			Committer: &protobuf.Contact{
				Name:  c.Committer.Name,
				Email: c.Committer.Email,
			},
			Parents:    parents,
			CommitTime: ts.New(c.Committer.TimeStamp),
		})
	}

	references := make([]*protobuf.Reference, len(refs))
	for i, rf := range refs {
		ref := protobuf.Reference{
			Name:      rf.Name,
			Shorthand: rf.Name,
			Hash:      string(rf.Target),
			IsRemote:  plumbing.IsRemote(rf.Name),
			IsBranch:  plumbing.IsBranch(rf.Name),
		}
		references[i] = &ref
	}

	headRef, _ := r.HEAD()
	head := protobuf.Head{
		Hash:      r.Peel(headRef).String(),
		Name:      headRef.Name,
		Shorthand: headRef.Name,
	}

	_ = protobuf.Repository{
		Path:       "",
		Commits:    commits,
		References: references,
		Head:       &head,
	}

	fmt.Println(len(commits))
	log.Printf("Log all commits took %s", time.Since(start))
}

func main() {
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:3000", nil))
	// }()

	startService()
}
