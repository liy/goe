package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-git/go-git/v5"
	goplumbing "github.com/go-git/go-git/v5/plumbing"
	goe "github.com/liy/goe/git"
	"github.com/liy/goe/plumbing"
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
	// cIter, err := r.Log(&git.LogOptions{All: true})
	// if err != nil {
	// 	return err
	// }

	// var commits []*protobuf.Commit
	// err = cIter.ForEach(func(c *object.Commit) error {
	// 	messages := strings.Split(c.Message, "\n")

	// 	summary := messages[0]
	// 	body := ""
	// 	if len(messages) > 1 {
	// 		body = strings.Join(messages[1:], "\n")
	// 	}

	// 	parents := make([]string, c.NumParents())
	// 	for i, pc := range c.ParentHashes {
	// 		parents[i] = pc.String()
	// 	}

	// 	commit := protobuf.Commit{
	// 		Hash:    c.Hash.String(),
	// 		Summary: summary,
	// 		Body:    body,
	// 		Author: &protobuf.Contact{
	// 			Name:  c.Author.Name,
	// 			Email: c.Author.Email,
	// 		},
	// 		Committer: &protobuf.Contact{
	// 			Name:  c.Committer.Name,
	// 			Email: c.Committer.Email,
	// 		},
	// 		Parents:    parents,
	// 		CommitTime: ts.New(c.Committer.When),
	// 	}
	// 	commits = append(commits, &commit)

	// 	return nil
	// })
	c, err := r.CommitObject(goplumbing.NewHash("4f3b0254d9160fd8786d2edb3a6a73ffcf6b70ac"))
	if err != nil {
		return err
	}
	fmt.Println(c)
	// _, err = r.CommitObject(goplumbing.NewHash("8e0228daabdc4708fb3f333fb869de84d5ed7d01"))
	// if err != nil {
	// 	return err
	// }
	// _, err = r.CommitObject(goplumbing.NewHash("ff46c79f1154922d155dcd7b1d18027ab265b2fa"))
	// if err != nil {
	// 	return err
	// }
	// _, err = r.CommitObject(goplumbing.NewHash("7d9095383a9a222fa3ba82eb8a803bcb338ad946"))
	// if err != nil {
	// 	return err
	// }
	// _, err = r.CommitObject(goplumbing.NewHash("77ace532f338d733006de2a34783201ca9d2bcc8"))
	// if err != nil {
	// 	return err
	// }
	// _, err = r.CommitObject(goplumbing.NewHash("1d90bf4d8d24d47e5fb3ac07aabf5242c96a6c31"))
	// if err != nil {
	// 	return err
	// }

	// fmt.Println(commit)

	// var references []*protobuf.Reference
	// rIter, err := r.References()
	// if err != nil {
	// 	return err
	// }

	// err = rIter.ForEach(func(r *goplumbing.Reference) error {
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


	// start := time.Now()
	// packReader := packfile.NewPackReader(largePack)

	// object, err := packReader.ReadObject(plumbing.ToHash("f9a08e80a692542cb94be651a61f81dd7374b39f"))
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(object)
	
	

	start := time.Now()
	r, err := goe.OpenRepository("./repo-test")
	if err != nil {
		fmt.Println(err)
	}
	o, err := r.GetTag(plumbing.ToHash("6dc409404e870d70c38a5ce9554c359a4ff339ee"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(o)
	log.Printf("Operation took %s", time.Since(start))
	
	// testRepository()
}