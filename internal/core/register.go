package core

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Register interface {
	// Grpc gRPC服务
	Grpc(server *grpc.Server)

	// Gateway 网关服务
	Gateway(mux *runtime.ServeMux, opts *[]grpc.DialOption) (ctx context.Context, handlers Handlers)
}
