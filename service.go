package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/liy/goe/git"
	"github.com/liy/goe/object"
	"github.com/liy/goe/plumbing"
	"github.com/liy/goe/src/protobuf"
	"google.golang.org/grpc"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

func startService() {
	lis, err := net.Listen("tcp", ":18888")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	protobuf.RegisterRepositoryServiceServer(s, new(RepositoryService))
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type RepositoryService struct {}

func (service *RepositoryService) GetHead(ctx context.Context, req *protobuf.GetHeadRequest) (*protobuf.GetHeadResponse, error) {
	start := time.Now()
    defer log.Printf("get head took %s", time.Since(start))

	// Opens an already existing repository.
	r, err := git.SimpleOpen(req.Path)
	if err != nil {
		return nil, fmt.Errorf("cannot open repository: %v", req.Path)
	}
	
	// ... retrieves the branch pointed by HEAD
	ref, err := r.HEAD()
	head := protobuf.Head {
		Hash: r.Peel(ref).String(),
		Name: ref.Name,
		Shorthand: ref.Shorthand(),
	}
	if err != nil {
		return nil, err
	}

	return &protobuf.GetHeadResponse{Head: &head}, nil
}

func (service *RepositoryService) GetRepository(ctx context.Context, req *protobuf.GetRepositoryRequest) (*protobuf.GetRepositoryResponse, error) {
	start := time.Now()
    defer log.Printf("log all commits took %s", time.Since(start))

	// Opens an already existing repository.
	r, err := git.SimpleOpen(req.Path)
	if err != nil {
		return nil, fmt.Errorf("cannot open repository: %v", req.Path)
	}
	
	refs := r.GetReferences()

	// annotated tags
	tags := make([]*protobuf.Tag, 0)

	// Setup potential tips
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

			tags = append(tags, &protobuf.Tag{
				Hash: tag.Hash.String(),
				Name: tag.Name,
				Message: tag.Message,
				Tagger: &protobuf.Signature{
					Name: tag.Tagger.Name,
					Email: tag.Tagger.Email,
				},
				Target: tag.Target.String(),
			})

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
			Author: &protobuf.Signature{
				Name:  strings.ToValidUTF8(c.Author.Name, ""),
				Email: c.Author.Email,
			},
			Committer: &protobuf.Signature{
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
			Shorthand: rf.Shorthand(),
			Hash:      string(rf.Target),
			IsRemote:  plumbing.IsRemote(rf.Name),
			IsBranch:  plumbing.IsBranch(rf.Name),
			IsTag:     plumbing.IsTag(rf.Name),
		}
		references[i] = &ref
	}

	headRef, _ := r.HEAD()
	head := protobuf.Head{
		Hash:      r.Peel(headRef).String(),
		Name:      headRef.Name,
		Shorthand: headRef.Shorthand(),
	}

	repository := protobuf.Repository{
		Path:       req.Path,
		Commits:    commits,
		References: references,
		Head:       &head,
		Tags: 	 tags,
	}

    log.Printf("commits %v", len(commits))

	return &protobuf.GetRepositoryResponse{Repository: &repository}, nil
}