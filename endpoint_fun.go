package grpc

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// EndpointHandler 端点注册方法
type EndpointHandler func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error
