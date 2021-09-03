package main

import (
	"fmt"
	"io/ioutil"

	"github.com/liy/goe/plumbing/indexfile"
	"github.com/liy/goe/plumbing/packfile"
)

func main() {
	// const port = ":8888"
	// listener, err := net.Listen("tcp", port)
	// CheckIfError(err)
	// credentials, err := credentials.NewServerTLSFromFile("./certificates/server.pem", "./certificates/server.key")
	// CheckIfError(err)
	// opts := []grpc.ServerOption{grpc.Creds(credentials), grpc.MaxRecvMsgSize(20 * 1024 * 1024), grpc.MaxSendMsgSize(20 * 1024 * 1024),}
	// s := grpc.NewServer(opts...)
	// protobuf.RegisterRepositoryServiceServer(s, new(RepositoryService))
	// s.Serve(listener)

	// CheckIfError(err)

	bytes, _ := ioutil.ReadFile("C:\\Source\\git-internals\\.git\\objects\\pack\\pack-66916c151da20048086dacbba45c420c0c1de8f6.idx")
	indexFile := new(indexfile.Index)
	indexFile.Decode(bytes)

	bytes, _ = ioutil.ReadFile("C:\\Source\\git-internals\\.git\\objects\\pack\\pack-66916c151da20048086dacbba45c420c0c1de8f6.pack")
	packFile := new(packfile.Pack)
	err := packFile.Decode(bytes)
	if err != nil {
		fmt.Println(err)
	}
}
