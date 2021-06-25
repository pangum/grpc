package grpc

import (
	`net`
	`time`

	`github.com/storezhang/pangu`
	`google.golang.org/grpc`
	`google.golang.org/grpc/keepalive`
)

const (
	initialWindowSize     = 1 << 30
	initialConnWindowSize = 1 << 30
	maxSendMsgSize        = 4 << 30
	maxRecvMsgSize        = 4 << 30
	keepAliveTime         = time.Second * 10
	keepAliveTimeout      = time.Second * 3
)

type gRPC struct {
	server   *grpc.Server
	services []func()
}

func New(options ...grpc.ServerOption) *gRPC {
	defaults := defaultOptions()
	options = append(options, defaults...)

	return &gRPC{
		server: grpc.NewServer(options...),
	}
}

func (g *gRPC) AddService(services ...func()) {
	g.services = append(g.services, services...)
}

func (g *gRPC) Run(config *pangu.Config) (err error) {
	panguConfig := new(panguConfig)
	if err = config.Load(panguConfig); nil != err {
		return
	}

	var listener net.Listener
	if listener, err = net.Listen("tcp", panguConfig.GRPC.Server.Addr()); nil != err {
		return
	}

	for _, s := range g.services {
		s()
	}

	return g.server.Serve(listener)
}

func defaultOptions() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.InitialWindowSize(initialWindowSize),
		grpc.InitialConnWindowSize(initialConnWindowSize),
		grpc.MaxSendMsgSize(maxSendMsgSize),
		grpc.MaxRecvMsgSize(maxRecvMsgSize),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			PermitWithoutStream: true,
		}),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:              keepAliveTime,
			Timeout:           keepAliveTimeout,
			MaxConnectionIdle: keepAliveTimeout,
		}),
	}
}
