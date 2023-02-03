package grpc

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// Register 空白注册器
type Register struct{}

func (r *Register) Grpc(_ *grpc.Server) {}

func (r *Register) Gateway(_ *runtime.ServeMux, _ string, _ ...grpc.DialOption) (err error) {
	return
}

type register interface {
	// Grpc 注册gRPC服务
	Grpc(server *grpc.Server)

	// Gateway 网关转换服务
	Gateway(mux *runtime.ServeMux, endpoint string, opts ...grpc.DialOption) (err error)
}
