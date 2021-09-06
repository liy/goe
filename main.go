package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/liy/goe/plumbing/indexfile"
	"github.com/liy/goe/plumbing/packfile"
	"github.com/liy/goe/src/protobuf"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

const defaultPack = ".\\repo\\.git\\objects\\pack\\pack-004ad14387e8ad228175d6e87e3281f0bd6b4d7e.pack"
const defaultPackIndex = ".\\repo\\.git\\objects\\pack\\pack-004ad14387e8ad228175d6e87e3281f0bd6b4d7e.idx"

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
	ref, err := r.Head()
	head := protobuf.Head{
		Hash:      ref.Hash().String(),
		Name:      ref.Name().String(),
		Shorthand: ref.Name().Short(),
	}
	if err != nil {
		return err
	}

	// ... retrieves the commit history
	cIter, err := r.Log(&git.LogOptions{All: true})
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
	if err != nil {
		return err
	}

	var references []*protobuf.Reference
	rIter, err := r.References()
	if err != nil {
		return err
	}

	err = rIter.ForEach(func(r *plumbing.Reference) error {
		ref := protobuf.Reference{
			Name:      r.Name().String(),
			Shorthand: r.Name().Short(),
			Hash:      r.Hash().String(),
			Type:      protobuf.Reference_Type(r.Type()),
			IsRemote:  r.Target().IsRemote(),
			IsBranch:  r.Target().IsBranch(),
		}
		references = append(references, &ref)
		return nil
	})
	if err != nil {
		return err
	}

	repository = protobuf.Repository{
		Path:       path,
		Commits:    commits,
		References: references,
		Head:       &head,
	}
    log.Printf("Log all commits took %s", time.Since(start))

	return nil
}

func main() {
	// testRepository()

	// const port = ":8888"
	// listener, err := net.Listen("tcp", port)
	// CheckIfError(err)
	// credentials, err := credentials.NewServerTLSFromFile("./certificates/server.pem", "./certificates/server.key")
	// CheckIfError(err)
	// opts := []grpc.ServerOption{grpc.Creds(credentials), grpc.MaxRecvMsgSize(20 * 1024 * 1024), grpc.MaxSendMsgSize(20 * 1024 * 1024)}
	// s := grpc.NewServer(opts...)
	// protobuf.RegisterRepositoryServiceServer(s, new(RepositoryService))
	// s.Serve(listener)
	// CheckIfError(err)

	start := time.Now()
	bytes, _ := ioutil.ReadFile(scratchPackIndex)
	indexFile := new(indexfile.Index)
	indexFile.Decode(bytes)

	bytes, _ = ioutil.ReadFile(scratchPack)
	packFile := new(packfile.Pack)
	// err := packFile.DecodeWithIndex(bytes, indexFile)
	err := packFile.Decode(bytes)
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("Log all commits took %s", time.Since(start))
}
