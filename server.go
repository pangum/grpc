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
	logger logging.Logger
}

func newServer(config *pangu.Config, logger logging.Logger) (server *Server, err error) {
	wrap := new(wrapper)
	if err = config.Load(wrap); nil != err {
		return
	}

	// 组织配置项
	conf := wrap.Grpc
	opts := make([]grpc.ServerOption, 0, 8)
	opts = append(opts, grpc.InitialWindowSize(int32(conf.Options.Size.Window.Initial)))
	opts = append(opts, grpc.InitialConnWindowSize(int32(conf.Options.Size.Window.Connection)))
	opts = append(opts, grpc.MaxSendMsgSize(int(conf.Options.Size.Msg.Send)))
	opts = append(opts, grpc.MaxRecvMsgSize(int(conf.Options.Size.Msg.Receive)))
	opts = append(opts, grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		PermitWithoutStream: conf.Options.Keepalive.Policy.Permit,
	}))
	opts = append(opts, grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: conf.Options.Keepalive.Idle,
		Time:              conf.Options.Keepalive.Time,
		Timeout:           conf.Options.Keepalive.Timeout,
	}))

	server = &Server{
		rpc:    grpc.NewServer(opts...),
		mux:    http.NewServeMux(),
		config: conf.Server,

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

	// 启动服务
	err = s.startup(listener)

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

func (s *Server) grpc(register register) (err error) {
	register.Grpc(s.rpc)

	return
}
