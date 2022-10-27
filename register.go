package grpc

import (
	"google.golang.org/grpc"
)

type register interface {
	// Grpc 注册gRPC服务
	Grpc(server *grpc.Server)
}
