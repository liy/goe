package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	goe "github.com/liy/goe/git"
	goeObject "github.com/liy/goe/object"
	"github.com/liy/goe/plumbing"
	"github.com/liy/goe/src/protobuf"
	ts "google.golang.org/protobuf/types/known/timestamppb"

	_ "net/http/pprof"

	_ "github.com/pkg/profile"
)

const defaultPack = ".\\repo-test\\.git\\objects\\pack\\pack-66916c151da20048086dacbba45c420c0c1de8f6.pack"
const defaultPackIndex = ".\\repo-test\\.git\\objects\\pack\\pack-66916c151da20048086dacbba45c420c0c1de8f6.idx"

const largePack = ".\\repo\\.git\\objects\\pack\\pack-004ad14387e8ad228175d6e87e3281f0bd6b4d7e.pack"
const largePackIndex = ".\\repo\\.git\\objects\\pack\\pack-004ad14387e8ad228175d6e87e3281f0bd6b4d7e.idx"

const scratchPack = "C://Users//liy//Workspace//goe//repo-scratch//.git//objects//pack//pack-1a1312067da0fc58cabf0a667c5fa43924928181.pack"
const scratchPackIndex = "C://Users//liy//Workspace//goe//repo-scratch//.git//objects//pack//pack-1a1312067da0fc58cabf0a667c5fa43924928181.idx"

func testRepository() error {
	path := "./repo"

	start := time.Now()

	// Opens an already existing repository.
	r, err := git.PlainOpen(path)
	if err != nil {
		return fmt.Errorf("cannot open repository: %v", path)
	}

	// ... retrieves the branch pointed by HEAD
	// ref, err := r.Head()
	// head := protobuf.Head{
	// 	Hash:      ref.Hash().String(),
	// 	Name:      ref.Name().String(),
	// 	Shorthand: ref.Name().Short(),
	// }
	// if err != nil {
	// 	return err
	// }

	// ... retrieves the commit history
	cIter, err := r.Log(&git.LogOptions{All: true, Order: git.LogOrderDefault})
	if err != nil {
		return err
	}

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

		commit := protobuf.Commit{
			Hash:    c.Hash.String(),
			Summary: summary,
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
			CommitTime: ts.New(c.Committer.When),
		}
		commits = append(commits, &commit)

		return nil
	})

	// var references []*protobuf.Reference
	// rIter, err := r.References()
	// if err != nil {
	// 	return err
	// }

	// err = rIter.ForEach(func(r *plumbing.Reference) error {
	// 	ref := protobuf.Reference{
	// 		Name:      r.Name().String(),
	// 		Shorthand: r.Name().Short(),
	// 		Hash:      r.Hash().String(),
	// 		Type:      protobuf.Reference_Type(r.Type()),
	// 		IsRemote:  r.Target().IsRemote(),
	// 		IsBranch:  r.Target().IsBranch(),
	// 	}
	// 	references = append(references, &ref)
	// 	return nil
	// })
	// if err != nil {
	// 	return err
	// }

	// repository = protobuf.Repository{
	// 	Path:       path,
	// 	Commits:    commits,
	// 	References: references,
	// 	Head:       &head,
	// }
	fmt.Println(len(commits))
	log.Printf("Log all commits took %s", time.Since(start))

	return nil
}

func mine() {
	start := time.Now()

	r, err := goe.OpenRepository("./repos/rails")
	if err != nil {
		fmt.Println(err)
	}

	refs := r.GetReferences()

	// Setup potential tips
	tips := make([]*goeObject.Commit, len(refs))
	for i, ref := range refs {
		var c *goeObject.Commit

		hashStr := r.TryPeel(ref.Target).String()
		raw, err := r.ReadObject(plumbing.ToHash(hashStr))
		if err != nil {
			fmt.Println(err)
			continue
		}

		if raw.Type == plumbing.OBJ_TAG {
			tag, err := goeObject.DecodeTag(raw)
			if err != nil {
				fmt.Println(err)
				continue
			}

			raw, err = r.ReadObject(tag.Target)
			if err != nil {
				continue
			}
		}

		c, err = goeObject.DecodeCommit(raw)
		if err != nil {
			fmt.Println(err)
			continue
		}

		tips[i] = c
	}

	var commits []*protobuf.Commit
	cItr := goe.NewCommitIterator(r, tips)
	for {
		c, err := cItr.Next()
		if err == goe.Done {
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
		Hash:      r.TryPeel(headRef.Target).String(),
		Name:      headRef.Name,
		Shorthand: headRef.Name,
	}

	_ = protobuf.Repository{
		Path:       r.GetPath(),
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

	// const port = ":8888"
	// listener, err := net.Listen("tcp", port)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// credentials, err := credentials.NewServerTLSFromFile("./certificates/server.pem", "./certificates/server.key")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// opts := []grpc.ServerOption{grpc.Creds(credentials), grpc.MaxRecvMsgSize(20 * 1024 * 1024), grpc.MaxSendMsgSize(20 * 1024 * 1024)}
	// s := grpc.NewServer(opts...)
	// protobuf.RegisterRepositoryServiceServer(s, new(RepositoryService))
	// s.Serve(listener)

	// testRepository()
	mine()
	// processed()
	// defer profile.Start().Stop()

	// r, err := goe.OpenRepository("./repo")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// raw, err := r.ReadObject(plumbing.ToHash("93e63b972b745f8766f5e529bc346e0d411ae36e"))
	// fmt.Println(raw.Type)

	// fmt.Println(r.GetReferences())
}
