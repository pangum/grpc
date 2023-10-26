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

	if rpc, gateway, le := s.listeners(); nil != le {
		err = le
	} else if gre := s.grpc(register, rpc); nil != gre {
		err = gre
	} else if gwe := s.gateway(register); nil != gwe {
		err = gwe
	} else {
		err = s.startup(gateway)
	}

	return
}

func (s *Server) Stop() (err error) {
	s.rpc.GracefulStop()
	err = s.http.Shutdown(context.Background())

	return
}

func (s *Server) startup(listener net.Listener) (err error) {
	s.http = new(http.Server)
	s.http.Addr = s.config.Gateway.Addr()
	s.http.Handler = s.handler(s.rpc, s.mux)
	s.http.ReadTimeout = s.config.Gateway.Timeout.Read
	s.http.ReadHeaderTimeout = s.config.Gateway.Timeout.Header

	fields := gox.Fields[any]{
		field.New("port.grpc", s.config.Server.Port),
	}
	if s.config.Gateway.Enable() {
		fields = append(fields, field.New("port.gateway", s.config.Gateway.Port))
	}
	s.logger.Info("启动gRPC服务器", fields...)
	err = s.http.Serve(listener)

	return
}

func (s *Server) grpc(register Register, listener net.Listener) (err error) {
	register.Grpc(s.rpc)
	if nil == s.config.Gateway || (s.config.Gateway.Enable() && s.diff()) {
		go s.serveRpc(listener)
	}

	return
}

func (s *Server) serveRpc(listener net.Listener) {
	if err := s.rpc.Serve(listener); nil != err {
		s.logger.Error("启动gRPC出错", field.New("addr", s.config.Server.Addr()))
	}
}

func (s *Server) diff() bool {
	return s.config.Gateway.Port != s.config.Server.Port
}

func (s *Server) listeners() (rpc net.Listener, gateway net.Listener, err error) {
	if listener, re := net.Listen(constant.Tcp, s.config.Server.Addr()); nil != re { // gRPC端口必须监听
		err = re
	} else if nil != s.config.Gateway && s.config.Gateway.Enable() && s.diff() { // 如果网关开启且端口不一样
		rpc = listener
		gateway, err = net.Listen(constant.Tcp, s.config.Gateway.Addr())
	} else { // 其它情况，监听端口都是一样的
		rpc = listener
		gateway = listener
	}

	return
}
