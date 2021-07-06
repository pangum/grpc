package grpc

import (
	`net`

	`github.com/storezhang/pangu`
	`google.golang.org/grpc`
	`google.golang.org/grpc/keepalive`
)

// Server gRPC服务器封装
type Server struct {
	grpc   *grpc.Server
	config serverConfig
}

func newServer(config *pangu.Config) (server *Server, err error) {
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
	}

	return
}

func (s *Server) Serve(functions ...registerFunc) (err error) {
	var listener net.Listener
	if listener, err = net.Listen("tcp", s.config.Addr()); nil != err {
		return
	}

	// 注册服务
	for _, function := range functions {
		function(s)
	}
	err = s.grpc.Serve(listener)

	return
}

func (s *Server) Server() *grpc.Server {
	return s.grpc
}
