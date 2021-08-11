package main

import (
	"net"

	"github.com/liy/goe/src/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	const port = ":8888"
	listener, err := net.Listen("tcp", port)
	CheckIfError(err)
	credentials, err := credentials.NewServerTLSFromFile("./certificates/server.pem", "./certificates/server.key")
	CheckIfError(err)
	opts := []grpc.ServerOption{grpc.Creds(credentials), grpc.MaxRecvMsgSize(20 * 1024 * 1024), grpc.MaxSendMsgSize(20 * 1024 * 1024),}
	s := grpc.NewServer(opts...)
	protobuf.RegisterRepositoryServiceServer(s, new(RepositoryService))
	s.Serve(listener)
	
	CheckIfError(err)
}