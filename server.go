package grpc

import (
	`net`

	`github.com/storezhang/glog`
	`github.com/storezhang/gox/field`
	`github.com/storezhang/pangu`
	`google.golang.org/grpc`
	`google.golang.org/grpc/keepalive`
	`google.golang.org/grpc/reflection`
)

// Server gRPC服务器封装
type Server struct {
	grpc   *grpc.Server
	config serverConfig

	logger glog.Logger
}

func newServer(config *pangu.Config, logger glog.Logger) (server *Server, err error) {
	panguConfig := new(panguConfig)
	if err = config.Load(panguConfig); nil != err {
		return
	}

	// 组织配置项
	grpcConfig := panguConfig.Grpc
	options := make([]grpc.ServerOption, 0, 8)
	options = append(options, grpc.InitialWindowSize(grpcConfig.Options.Size.Window.Initial))
	options = append(options, grpc.InitialConnWindowSize(grpcConfig.Options.Size.Window.Connection))
	options = append(options, grpc.MaxSendMsgSize(grpcConfig.Options.Size.Msg.Send))
	options = append(options, grpc.MaxRecvMsgSize(grpcConfig.Options.Size.Msg.Receive))
	options = append(options, grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		PermitWithoutStream: grpcConfig.Options.Keepalive.Policy.Permit,
	}))
	options = append(options, grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: grpcConfig.Options.Keepalive.Idle,
		Time:              grpcConfig.Options.Keepalive.Time,
		Timeout:           grpcConfig.Options.Keepalive.Timeout,
	}))

	server = &Server{
		grpc:   grpc.NewServer(options...),
		config: grpcConfig.Server,

		logger: logger,
	}

	return
}

func (s *Server) Serve(fun registerFunc, opts ...serveOption) (err error) {
	options := defaultServeOptions()
	for _, opt := range opts {
		opt.apply(options)
	}

	var listener net.Listener
	if listener, err = net.Listen("tcp", s.config.Addr()); nil != err {
		return
	}

	// 注册服务
	fun()

	// 处理选项
	if options.Reflection {
		reflection.Register(s.grpc)
	}

	s.logger.Info("启动gRPC服务器", field.Int("port", s.config.Port))
	// 启动服务
	err = s.grpc.Serve(listener)

	return
}

func (s *Server) Grpc() *grpc.Server {
	return s.grpc
}

func (s *Server) Stop() {
	s.grpc.GracefulStop()
}
