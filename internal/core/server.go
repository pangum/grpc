package core

import (
	"context"
	"net"
	"net/http"

	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/goexl/log"
	"github.com/pangum/grpc/internal/config"
	"github.com/pangum/grpc/internal/core/internal"
	"github.com/pangum/grpc/internal/internal/constant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server gRPC服务器封装
type Server struct {
	rpc  *grpc.Server
	http *http.Server
	mux  *http.ServeMux

	config *internal.Config
	logger log.Logger
	_      gox.CannotCopy
}

func NewServer(
	rpc *grpc.Server,
	mux *http.ServeMux,
	server *config.Server, gateway *config.Gateway,
	logger log.Logger,
) *Server {
	return &Server{
		rpc: rpc,
		mux: mux,

		config: internal.NewConfig(server, gateway),
		logger: logger,
	}
}

func (s *Server) Serve(register Register) (err error) {
	if *s.config.Server.Reflection { // 反射，在gRPC接口调试时，可以反射出方法和参数
		reflection.Register(s.rpc)
	}

	if listener, le := net.Listen(constant.Tcp, s.config.Server.Addr()); nil != le {
		err = le
	} else if gre := s.grpc(register); nil != gre {
		err = gre
	} else if gwe := s.gateway(register); nil != gwe {
		err = gwe
	} else {
		err = s.startup(listener)
	}

	return
}

func (s *Server) Stop() {
	s.rpc.GracefulStop()
	_ = s.http.Shutdown(context.Background())
}

func (s *Server) startup(listener net.Listener) (err error) {
	s.http = new(http.Server)
	s.http.Addr = s.config.Gateway.Addr()
	s.http.Handler = s.handler(s.rpc, s.mux)
	s.http.ReadTimeout = s.config.Gateway.Timeout.Read
	s.http.ReadHeaderTimeout = s.config.Gateway.Timeout.Header

	fields := gox.Fields[any]{
		field.New("port.rpc", s.config.Server.Port),
	}
	if s.config.Gateway.Enable() {
		fields = append(fields, field.New("port.gateway", s.config.Gateway.Port))
	}
	s.logger.Info("启动gRPC服务器", fields...)
	err = s.http.Serve(listener)

	return
}

func (s *Server) grpc(register Register) (err error) {
	register.Grpc(s.rpc)

	return
}
