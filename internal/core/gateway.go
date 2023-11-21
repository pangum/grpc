package core

import (
	"context"
	"net/http"
	"strconv"

	"github.com/goexl/exc"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/goexl/log"
	"github.com/pangum/grpc/internal/internal/constant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Gateway struct {
	logger log.Logger
	_      gox.CannotCopy
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

func (g *Gateway) Error(code codes.Code, err error) error {
	return status.Error(code, err.Error())
}

func (g *Gateway) Exception(code codes.Code, exception exc.Exception) error {
	return status.Error(code, exception.Error())
}

func (g *Gateway) ValidationError(err error) error {
	return g.Error(http.StatusBadRequest, err)
}

func (g *Gateway) ServerError(err error) error {
	return g.Error(http.StatusInternalServerError, err)
}

func (g *Gateway) ServerException(code int, fields gox.Fields[any]) (err error) {
	exception := exc.NewException(code, "服务器错误，客户端需要根据返回中的`code`码来确认具体是什么错误", fields...)
	err = g.Exception(http.StatusInternalServerError, exception)

	return
}

func (g *Gateway) Notfound() error {
	return g.Error(http.StatusNotFound, exc.NewMessage("未找到资源"))
}
