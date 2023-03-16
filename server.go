package grpc

import (
	"context"
	"net"
	"net/http"

	"github.com/goexl/gox/field"
	"github.com/pangum/logging"
	"github.com/pangum/pangu"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// Server gRPC服务器封装
type Server struct {
	rpc    *grpc.Server
	http   *http.Server
	mux    *http.ServeMux
	config server

	serveHttp bool
	logger    *logging.Logger
}

func newServer(config *pangu.Config, logger *logging.Logger) (server *Server, err error) {
	_panguConfig := new(panguConfig)
	if err = config.Load(_panguConfig); nil != err {
		return
	}

	// 组织配置项
	_config := _panguConfig.Grpc
	_options := make([]grpc.ServerOption, 0, 8)
	_options = append(_options, grpc.InitialWindowSize(int32(_config.Options.Size.Window.Initial)))
	_options = append(_options, grpc.InitialConnWindowSize(int32(_config.Options.Size.Window.Connection)))
	_options = append(_options, grpc.MaxSendMsgSize(int(_config.Options.Size.Msg.Send)))
	_options = append(_options, grpc.MaxRecvMsgSize(int(_config.Options.Size.Msg.Receive)))
	_options = append(_options, grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		PermitWithoutStream: _config.Options.Keepalive.Policy.Permit,
	}))
	_options = append(_options, grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: _config.Options.Keepalive.Idle,
		Time:              _config.Options.Keepalive.Time,
		Timeout:           _config.Options.Keepalive.Timeout,
	}))

	server = &Server{
		rpc:    grpc.NewServer(_options...),
		mux:    http.NewServeMux(),
		config: _config.Server,

		logger: logger,
	}

	return
}

func (s *Server) Serve(register register, opts ...serveOption) (err error) {
	_options := defaultServeOptions()
	for _, opt := range opts {
		opt.apply(_options)
	}

	var listener net.Listener
	if listener, err = net.Listen("tcp", s.config.Addr()); nil != err {
		return
	}

	// 注册服务
	if err = s.grpc(register); nil != err {
		return
	}
	if err = s.gateway(register); nil != err {
		return
	}
	if err = s.metric(register); nil != err {
		return
	}

	// 处理选项
	if _options.Reflection {
		reflection.Register(s.rpc)
	}

	s.logger.Info("启动gRPC服务器", field.New("port", s.config.Port))
	// 启动服务
	err = s.startup(listener)

	return
}

func (s *Server) Stop() {
	if s.serveHttp {
		_ = s.http.Shutdown(context.Background())
	} else {
		s.rpc.GracefulStop()
	}
}

func (s *Server) startup(listener net.Listener) (err error) {
	if s.serveHttp {
		s.http = &http.Server{
			Addr:              s.config.Addr(),
			Handler:           s.handler(s.rpc, s.mux),
			ReadTimeout:       s.config.Timeout.Read,
			ReadHeaderTimeout: s.config.Timeout.Header,
		}
		err = s.http.Serve(listener)
	} else {
		err = s.rpc.Serve(listener)
	}

	return
}

func (s *Server) grpc(register register) (err error) {
	register.Grpc(s.rpc)

	return
}
