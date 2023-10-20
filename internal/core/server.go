package core

import (
	"context"
	"net"
	"net/http"

	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/goexl/simaqian"
	"github.com/pangum/grpc/internal"
	"github.com/pangum/grpc/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server gRPC服务器封装
type Server struct {
	rpc    *grpc.Server
	http   *http.Server
	mux    *http.ServeMux
	config *config.Server

	logger simaqian.Logger
	_      gox.CannotCopy
}

func NewServer(server *grpc.Server, mux *http.ServeMux, config *config.Server, logger simaqian.Logger) *Server {
	return &Server{
		rpc:    server,
		mux:    mux,
		config: config,

		logger: logger,
	}
}

func (s *Server) Serve(register Register) (err error) {
	if *s.config.Reflection { // 反射，在gRPC接口调试时，可以反射出方法和参数
		reflection.Register(s.rpc)
	}

	if listener, le := net.Listen(internal.Tcp, s.config.Addr()); nil != le {
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
	s.http.Addr = s.config.Addr()
	s.http.Handler = s.handler(s.rpc, s.mux)
	s.http.ReadTimeout = s.config.Timeout.Read
	s.http.ReadHeaderTimeout = s.config.Timeout.Header

	s.logger.Info("启动gRPC服务器", field.New("port", s.config.Port))
	err = s.http.Serve(listener)

	return
}

func (s *Server) grpc(register Register) (err error) {
	register.Grpc(s.rpc)

	return
}
