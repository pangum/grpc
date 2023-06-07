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
	logging.Logger

	rpc    *grpc.Server
	http   *http.Server
	mux    *http.ServeMux
	config server
}

func newServer(config *pangu.Config, logger logging.Logger) (server *Server, mux *http.ServeMux, err error) {
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

	mux = http.NewServeMux()
	server = &Server{
		Logger: logger,

		rpc:    grpc.NewServer(opts...),
		mux:    mux,
		config: conf.Server,
	}

	return
}

func (s *Server) Serve(register register, opts ...serveOption) (err error) {
	_options := defaultServeOptions()
	for _, opt := range opts {
		opt.apply(_options)
	}

	if _options.Reflection { // 反射，在gRPC接口调试时，可以反射出方法和参数
		reflection.Register(s.rpc)
	}

	if listener, le := net.Listen(tcp, s.config.Addr()); nil != le {
		err = le
	} else if re := s.grpc(register); nil != re {
		err = re
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

	s.Info("启动gRPC服务器", field.New("port", s.config.Port))
	err = s.http.Serve(listener)

	return
}

func (s *Server) grpc(register register) (err error) {
	register.Grpc(s.rpc)

	return
}
