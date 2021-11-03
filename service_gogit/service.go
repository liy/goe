package service_gogit

import (
	"context"
	"fmt"
	"log"
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

func (service *RepositoryService) GetHead(ctx context.Context, req *protobuf.EmptyRequest) (*protobuf.GetHeadResponse, error) {
	md, ok :=  metadata.FromIncomingContext(ctx);
	if !ok {
		return nil, fmt.Errorf("cannot retrieve metadata");
	}

	path := "./repo"
	if mdValues := md.Get("path"); len(mdValues) != 0 {
		path = mdValues[0]
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

    log.Printf("Get head took %s", time.Since(start))

	return &protobuf.GetHeadResponse{Head: &head}, nil
}

func (service *RepositoryService) GetRepository(ctx context.Context, req *protobuf.EmptyRequest) (*protobuf.GetRepositoryResponse, error) {
	md, ok :=  metadata.FromIncomingContext(ctx);
	if !ok {
		return nil, fmt.Errorf("cannot retrieve metadata");
	}

	path := "./repo"
	if mdValues := md.Get("path"); len(mdValues) != 0 {
		path = mdValues[0]
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
	cIter, err := r.Log(&git.LogOptions{All: true})
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

		if(c.Hash.String() == "99bc896d3d914c1607b8ee99b9f2cb51e2fd2b28") {
			fmt.Println("!!!")
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

// func (service *RepositoryService) Watch(req *protobuf.EmptyRequest, stream protobuf.RepositoryService_WatchServer ) error {
// 	md, ok := metadata.FromIncomingContext(stream.Context()); 
// 	if !ok {
// 		return fmt.Errorf("cannot retrieve metadata");
// 	}

// 	// TODO: when .git folder changes check head


// 	path := "./repo"
// 	if mdValues := md.Get("path"); len(mdValues) != 0 {
// 		path = mdValues[0]
// 	}

// 	// Opens an already existing repository.
// 	r, err := git.PlainOpen(path)
// 	if err != nil {
// 		return fmt.Errorf("cannot open repository: %v", path)
// 	}

// 	// ... retrieves the branch pointed by HEAD
// 	ref, err := r.Head()
// 	head := protobuf.Head {
// 		Hash: ref.Hash().String(),
// 		Name: ref.Name().String(),
// 		Shorthand: ref.Name().Short(),
// 	}
// 	if err != nil {
// 		return err
// 	}

// 	stream.Send(&protobuf.WatchResponse{Head: &head})

// 	return nil
// }

