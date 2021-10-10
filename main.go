package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	goe "github.com/liy/goe/git"
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

	r, err := goe.OpenRepository("./repo")
	if err != nil {
		fmt.Println(err)
	}

	// commits, _:=r.GetCommits(goePlumbing.ToHash("29970a12c888ed57ae7b723551238375c70eac08"))
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

	refs, err := r.GetReferences()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(refs)
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
}