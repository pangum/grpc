package core

import (
	"context"
	"strconv"

	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/goexl/log"
	"github.com/pangum/grpc/internal/internal/constant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Gateway struct {
	logger log.Logger
	_      gox.Pointerized
}

func NewGateway(logger log.Logger) *Gateway {
	return &Gateway{
		logger: logger,
	}
}

func (g *Gateway) Status(ctx context.Context, code int) {
	header := metadata.Pairs(constant.HttpStatusHeader, strconv.Itoa(code))
	if err := grpc.SetHeader(ctx, header); nil != err {
		g.logger.Debug("设置状态码出错", field.New("code", code), field.Error(err))
	}
}
