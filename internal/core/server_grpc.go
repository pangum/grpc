package core

import (
	"net"

	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
)

func (s *Server) setupGrpc(register Register, listener net.Listener) (err error) {
	register.Grpc(s.rpc)
	fields := gox.Fields[any]{
		field.New("name", s.config.Server.Name),
		field.New("addr", s.config.Server.Addr()),
	}
	s.logger.Info("启动服务成功", fields...)
	if !s.gatewayEnabled() || (s.gatewayEnabled() && s.diffPort()) {
		s.rpcStarted = true
		s.wait.Add(1)
		go s.serveRpc(listener, &fields)
	}

	return
}

func (s *Server) serveRpc(listener net.Listener, fields *gox.Fields[any]) {
	defer s.wait.Done()

	if err := s.rpc.Serve(listener); nil != err {
		s.logger.Error("启动服务出错", fields.Add(field.Error(err))...)
	}
}
