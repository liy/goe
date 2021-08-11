package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/liy/goe/src/protobuf"
	"google.golang.org/grpc/metadata"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

var repository protobuf.Repository

type RepositoryService struct {}

func (service *RepositoryService) GetCommits(ctx context.Context, req *protobuf.EmptyRequest) (*protobuf.GetCommitsResponse, error) {
	md, ok :=  metadata.FromIncomingContext(ctx);
	if !ok {
		return nil, fmt.Errorf("cannot retrieve metadata");
	}

	path := "./repo"
	if mdValues := md.Get("path"); len(mdValues) != 0 {
		path = mdValues[0]
	}

	batchSize := 2000
	if mdValues := md.Get("batchSize"); len(mdValues) != 0 {
		if value, err := strconv.Atoi(md.Get("batchSize")[0]); err != nil {
			batchSize = value
		}
	}

	start := time.Now()

	// Opens an already existing repository.
	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open repository: %v", path)
	}
	
	// ... retrieves the branch pointed by HEAD
	ref, err := r.Head()
	if err != nil {
		return nil, err
	}

	// ... retrieves the commit history
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash(), All: true, Order: git.LogOrderCommitterTime})
	if err != nil {
		return nil, err
	}

	commits := make([]*protobuf.Commit, batchSize)
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
		
		return nil
	})
	if err != nil {
		return nil, err
	}

    log.Printf("Log all commits took %s", time.Since(start))

	return &protobuf.GetCommitsResponse{Commits: commits}, nil
}

func (service *RepositoryService) GetRepository(ctx context.Context, req *protobuf.EmptyRequest) (*protobuf.GetRepositoryResponse, error) {
	path := "./repo"
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		path = md.Get("path")[0]
	}

	start := time.Now()

	// Opens an already existing repository.
	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open repository: %v", path)
	}
	
	
	// ... retrieves the branch pointed by HEAD
	ref, err := r.Head()
	head := protobuf.Head {
		Hash: ref.Hash().String(),
		Name: ref.Name().String(),
		Shorthand: ref.Name().Short(),
	}
	if err != nil {
		return nil, err
	}

	// ... retrieves the commit history
	since := time.Time{}
	until := time.Date(2021, 9, 30, 0, 0, 0, 0, time.UTC)
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash(), Since: &since, Until: &until, Order: git.LogOrderCommitterTime})
	if err != nil {
		return nil, err
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


		return nil
	})
	if err != nil {
		return nil, err
	}

	var references []*protobuf.Reference
	rIter, err := r.References()
	if err != nil {
		return nil, err
	}

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
	if err != nil {
		return nil, err
	}


	repository = protobuf.Repository {
		Path: path,
		Commits: commits,
		References: references,
		Head: &head,
	}

    log.Printf("Log all commits took %s", time.Since(start))

	return &protobuf.GetRepositoryResponse{Repository: &repository}, nil
}
