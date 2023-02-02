package grpc

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/goexl/gox/field"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pangum/logging"
	"github.com/pangum/pangu"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// Server gRPC服务器封装
type Server struct {
	grpc    *grpc.Server
	gateway *http.Server
	config  server

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
	if err = s.setupGrpc(register); nil != err {
		return
	}
	if err = s.setupGateway(register); nil != err {
		return
	}

	// 处理选项
	if _options.Reflection {
		reflection.Register(s.grpc)
	}

	s.logger.Info("启动gRPC服务器", field.New("port", s.config.Port))
	// 启动服务
	err = s.startup(listener)

	return
}

func (s *Server) Stop() {
	if s.config.gateway() {
		_ = s.gateway.Shutdown(context.Background())
	} else {
		s.grpc.GracefulStop()
	}
}

func (s *Server) startup(listener net.Listener) (err error) {
	if s.config.gateway() {
		err = s.gateway.Serve(listener)
	} else {
		err = s.grpc.Serve(listener)
	}

	return
}

func (s *Server) setupGrpc(register register) (err error) {
	register.Grpc(s.grpc)

	return
}

func (s *Server) setupGateway(register register) (err error) {
	if !s.config.gateway() {
		return
	}

	gateway := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err = register.Gateway(gateway, s.config.Addr(), opts...); nil != err {
		return
	}

	mux := http.NewServeMux()
	mux.Handle("/", gateway)
	s.gateway = &http.Server{
		Addr:    s.config.Addr(),
		Handler: s.handler(s.grpc, mux),
	}

	return
}

func (s *Server) handler(grpc *grpc.Server, gateway http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpc.ServeHTTP(w, r)
		} else {
			gateway.ServeHTTP(w, r)
		}
	}), new(http2.Server))
}
