package core

import (
	"context"
	"net"
	"net/http"

	"github.com/goexl/gox"
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

	if rpc, gateway, le := s.listeners(); nil != le {
		err = le
	} else if gre := s.setupGrpc(register, rpc); nil != gre {
		err = gre
	} else if gwe := s.setupGateway(register, gateway); nil != gwe {
		err = gwe
	}

	return
}

func (s *Server) Stop() (err error) {
	s.rpc.GracefulStop()
	err = s.http.Shutdown(context.Background())

	return
}

func (s *Server) diff() bool {
	return s.config.Gateway.Port != s.config.Server.Port
}

func (s *Server) listeners() (rpc net.Listener, gateway net.Listener, err error) {
	if listener, re := net.Listen(constant.Tcp, s.config.Server.Addr()); nil != re { // gRPC端口必须监听
		err = re
	} else if s.gatewayEnabled() && s.diff() { // 如果网关开启且端口不一样
		rpc = listener
		gateway, err = net.Listen(constant.Tcp, s.config.Gateway.Addr())
	} else { // 其它情况，监听端口都是一样的
		rpc = listener
		gateway = listener
	}

	return
}

func (s *Server) gatewayEnabled() bool {
	return nil != s.config.Gateway && s.config.Gateway.Enable()
}
