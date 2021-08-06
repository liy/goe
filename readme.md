# Build protocol buffer

Compile protobuf and grpc related protobuf services

> protoc --proto_path=./protobuf --go-grpc_out=require_unimplemented_servers=false:./ --go_out=./ ./protobuf/messages.proto

Note that `require_unimplemented_servers` is added to disable the backward-compatibility behaviour: https://github.com/grpc/grpc-go/releases/tag/cmd%2Fprotoc-gen-go-grpc%2Fv1.0.0
