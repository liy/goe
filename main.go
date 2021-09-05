package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

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

	// bytes, _ := ioutil.ReadFile("C:\\Source\\git-internals\\.git\\objects\\pack\\pack-66916c151da20048086dacbba45c420c0c1de8f6.idx")
	// indexFile := new(indexfile.Index)
	// indexFile.Decode(bytes)

	
	start := time.Now()
	bytes, _ := ioutil.ReadFile(".\\repo\\.git\\objects\\pack\\pack-004ad14387e8ad228175d6e87e3281f0bd6b4d7e.pack")
	packFile := new(packfile.Pack)
	err := packFile.Decode(bytes)
	if err != nil {
		fmt.Println(err)
	}
    log.Printf("Log all commits took %s", time.Since(start))
}
