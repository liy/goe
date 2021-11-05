package main

import (
	"fmt"
	"log"
	"time"

	"github.com/liy/goe/git"
	"github.com/liy/goe/object"
	"github.com/liy/goe/plumbing"

	_ "net/http/pprof"

	_ "github.com/pkg/profile"
)

func main() {
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:3000", nil))
	// }()

	// startService()

	r, _ := git.SimpleOpen("../repos/git")

	refs := r.GetReferences()
	tips := make([]*object.Commit, 0)
	for _, ref := range refs {
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

		tips = append(tips, c)
	}

	start := time.Now()
	git.NewCommitIterator(r, tips)

	cItr := git.NewCommitIterator(r, tips)
	commits := make([]*object.Commit, cItr.Size)
	idx := 0 
	for {
		c, err := cItr.Next()
		if err == git.Done {
			break
		}
		if err != nil {
			break
		}
		commits[idx] = c
		idx++
	}
	fmt.Println(cItr.Size)
    log.Printf("Log all commits took %s", time.Since(start))
}
