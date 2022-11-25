package grpc

import (
	"net"

	"github.com/goexl/gox/field"
	"github.com/pangum/logging"
	"github.com/pangum/pangu"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// Server gRPC服务器封装
type Server struct {
	grpc   *grpc.Server
	config server

	logger *logging.Logger
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
		grpc:   grpc.NewServer(_options...),
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
	register.Grpc(s.grpc)

	// 处理选项
	if _options.Reflection {
		reflection.Register(s.grpc)
	}

	s.logger.Info("启动gRPC服务器", field.New("port", s.config.Port))
	// 启动服务
	err = s.grpc.Serve(listener)

	return
}

func (s *Server) Stop() {
	s.grpc.GracefulStop()
}
