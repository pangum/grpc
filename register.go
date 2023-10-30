package grpc

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pangum/grpc/internal/core"
	"google.golang.org/grpc"
)

var _ core.Register = (*Register)(nil)

// Register 空白注册器
type Register struct{}

func (r *Register) Grpc(_ *grpc.Server) {}

func (r *Register) Gateway(_ *runtime.ServeMux, _ *[]grpc.DialOption) (ctx context.Context, handlers core.Handlers) {
	return
}

func (r *Register) Before(_ context.Context) (err error) {
	return
}

func (r *Register) After(_ context.Context) (err error) {
	return
}
