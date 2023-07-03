package grpc

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

var _ register = (*Register)(nil)

// Register 空白注册器
type Register struct{}

func (r *Register) Grpc(_ *grpc.Server) {}

func (r *Register) Gateway(_ *runtime.ServeMux, _ *[]grpc.DialOption) (ctx context.Context, handlers []EndpointHandler) {
	return
}

type register interface {
	// Grpc gRPC服务
	Grpc(server *grpc.Server)

	// Gateway 网关服务
	Gateway(mux *runtime.ServeMux, opts *[]grpc.DialOption) (ctx context.Context, handlers []EndpointHandler)
}
