package grpc

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type register interface {
	// Grpc 注册gRPC服务
	Grpc(server *grpc.Server)

	// Gateway 网关转换服务
	Gateway(mux *runtime.ServeMux, endpoint string, opts ...grpc.DialOption) (err error)
}
