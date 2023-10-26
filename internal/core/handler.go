package core

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// Handler 端点注册方法
type Handler func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error
