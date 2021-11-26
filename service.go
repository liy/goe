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
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

func startService() {
	const port = ":8888"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
	}
	credentials, err := credentials.NewServerTLSFromFile("./certificates/server.pem", "./certificates/server.key")
	if err != nil {
		fmt.Println(err)
	}

	opts := []grpc.ServerOption{grpc.Creds(credentials), grpc.MaxRecvMsgSize(100 * 1024 * 1024), grpc.MaxSendMsgSize(100 * 1024 * 1024)}
	s := grpc.NewServer(opts...)
	protobuf.RegisterRepositoryServiceServer(s, new(RepositoryService))
	s.Serve(listener)
}

type RepositoryService struct {}

func (service *RepositoryService) GetHead(ctx context.Context, req *protobuf.EmptyRequest) (*protobuf.GetHeadResponse, error) {
	md, ok :=  metadata.FromIncomingContext(ctx);
	if !ok {
		return nil, fmt.Errorf("cannot retrieve metadata");
	}

	var path = "../repos/topo-sort/"
	if mdValues := md.Get("path"); len(mdValues) != 0 {
		path = mdValues[0]
	}

	start := time.Now()

	// Opens an already existing repository.
	r, err := git.SimpleOpen(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open repository: %v", path)
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

    log.Printf("Get head took %s", time.Since(start))

	return &protobuf.GetHeadResponse{Head: &head}, nil
}

func (service *RepositoryService) GetRepository(ctx context.Context, req *protobuf.EmptyRequest) (*protobuf.GetRepositoryResponse, error) {
	md, ok :=  metadata.FromIncomingContext(ctx);
	if !ok {
		return nil, fmt.Errorf("cannot retrieve metadata");
	}

	var path = "../repos/topo-sort/"
	if mdValues := md.Get("path"); len(mdValues) != 0 {
		path = mdValues[0]
	}

	start := time.Now()

	// Opens an already existing repository.
	r, err := git.SimpleOpen(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open repository: %v", path)
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
		Path:       path,
		Commits:    commits,
		References: references,
		Head:       &head,
		Tags: 	 tags,
	}

    log.Printf("commits %v", len(commits))
    log.Printf("Log all commits took %s", time.Since(start))

	return &protobuf.GetRepositoryResponse{Repository: &repository}, nil
}