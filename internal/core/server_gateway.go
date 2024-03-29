package core

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pangum/grpc/internal/decoder"
	"github.com/pangum/grpc/internal/internal/constant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func (s *Server) setupGateway(register Register, listener net.Listener) (err error) {
	if !s.gatewayEnabled() {
		return
	}

	gatewayOpts := s.config.Gateway.Options()
	gatewayOpts = append(gatewayOpts, runtime.WithForwardResponseOption(s.response))
	gatewayOpts = append(gatewayOpts, runtime.WithIncomingHeaderMatcher(s.in))
	gatewayOpts = append(gatewayOpts, runtime.WithOutgoingHeaderMatcher(s.out))
	gatewayOpts = append(gatewayOpts, runtime.WithMetadata(s.metadata))
	gatewayOpts = append(gatewayOpts, runtime.WithMetadata(s.metadata))
	gatewayOpts = append(gatewayOpts, runtime.WithErrorHandler(s.error))
	// 使用特定的解码器来处理原始数据
	gatewayOpts = append(gatewayOpts, runtime.WithMarshalerOption(constant.RawHeaderValue, decoder.NewRaw()))
	if nil != s.config.Gateway.Unescape {
		gatewayOpts = append(gatewayOpts, runtime.WithUnescapingMode(s.config.Gateway.Unescape.Mode))
	}

	gw := runtime.NewServeMux(gatewayOpts...)
	grpcOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	ctx, handlers := register.Gateway(gw, &grpcOpts)
	if s.rpcStarted && !s.diffPort() {
		grpcOpts = append(grpcOpts, grpc.WithBlock())
	}
	if nil == ctx {
		ctx = context.Background()
	}
	if connection, dce := grpc.DialContext(ctx, s.config.Server.Addr(), grpcOpts...); nil != dce {
		err = dce
	} else if ge := s.registerGateway(ctx, gw, connection, handlers); nil != ge {
		err = ge
	} else if "" == s.config.Gateway.Path {
		s.mux.Handle(constant.Slash, gw)
	} else {
		path := s.config.Gateway.Path
		s.mux.Handle(gox.StringBuilder(path, constant.Slash).String(), http.StripPrefix(path, gw))
	}
	if nil == err {
		err = s.serveGateway(listener)
	}

	return
}

func (s *Server) serveGateway(listener net.Listener) (err error) {
	s.http = new(http.Server)
	s.http.Addr = s.config.Gateway.Addr()
	s.http.Handler = s.handler(s.rpc, s.mux)
	s.http.ReadTimeout = s.config.Gateway.Timeout.Read
	s.http.ReadHeaderTimeout = s.config.Gateway.Timeout.Header

	fields := gox.Fields[any]{
		field.New("name", s.config.Gateway.Name),
		field.New("addr", s.config.Gateway.Addr()),
	}
	s.logger.Info("启动服务成功", fields...)
	err = s.http.Serve(listener)

	return
}

func (s *Server) registerGateway(
	ctx context.Context,
	mux *runtime.ServeMux, connection *grpc.ClientConn,
	handlers []Handler,
) (err error) {
	for _, handler := range handlers {
		if he := handler(ctx, mux, connection); nil != he {
			err = he
			s.logger.Warn("注册网关出错", field.New("func", handler), field.Error(he))
		}
		if nil != err {
			break
		}
	}

	return
}

func (s *Server) response(ctx context.Context, writer http.ResponseWriter, msg proto.Message) (err error) {
	// 注意，这儿的顺序不能乱，必须先写入头再写入状态码
	if se := s.header(ctx, writer, msg); nil != se {
		err = se
	} else if he := s.status(ctx, writer); nil != he {
		err = he
	}

	return
}

func (s *Server) status(ctx context.Context, writer http.ResponseWriter) (err error) {
	if md, ok := runtime.ServerMetadataFromContext(ctx); !ok {
		// 上下文无法转换
	} else if _status := md.HeaderMD.Get(constant.HttpStatusHeader); 0 == len(_status) {
		// 没有设置状态
	} else if code, ae := strconv.Atoi(_status[0]); nil != ae {
		err = ae
		s.logger.Warn("状态码被错误设置", field.New("value", _status))
	} else {
		md.HeaderMD.Delete(constant.HttpStatusHeader)
		writer.Header().Del(constant.GrpcStatusHeader)
		writer.WriteHeader(code)
	}

	return
}

func (s *Server) error(
	_ context.Context, _ *runtime.ServeMux, _ runtime.Marshaler,
	writer http.ResponseWriter, _ *http.Request,
	err error,
) {
	if _status, ok := status.FromError(err); ok {
		original := _status.Code()
		code := gox.Ift(codes.Unknown == _status.Code(), http.StatusBadGateway, int(_status.Code()))
		// ! 修复状态码，避免返回错误的状态吗
		code = gox.Ift(s.isHttpCode(code), code, runtime.HTTPStatusFromCode(original))
		writer.WriteHeader(code)
		bytes := []byte(_status.Message())
		_, _ = writer.Write(bytes)
	} else {
		writer.WriteHeader(http.StatusBadGateway)
		bytes := []byte(err.Error())
		_, _ = writer.Write(bytes)
	}
}

func (s *Server) header(ctx context.Context, writer http.ResponseWriter, _ proto.Message) (err error) {
	if md, ok := runtime.ServerMetadataFromContext(ctx); ok {
		s.resetHeader(md.HeaderMD, writer)
	}

	return
}

func (s *Server) resetHeader(header metadata.MD, writer http.ResponseWriter) {
	for key, value := range header {
		if constant.HttpStatusHeader == key { // 不处理设置状态码的逻辑，由状态码设置逻辑特殊处理
			continue
		}

		newKey := strings.ToLower(key)
		removal := false
		newKey, removal = s.config.Gateway.Header.TestRemove(newKey)

		if removal {
			writer.Header().Set(newKey, strings.Join(value, constant.Space))
			header.Delete(key)
			writer.Header().Del(fmt.Sprintf(constant.GrpcMetadataFormatter, key))
		}
	}
}

func (s *Server) in(key string) (new string, match bool) {
	if newKey, test := s.config.Gateway.Header.TestIns(key); test {
		new = newKey
		match = true
	} else {
		new, match = runtime.DefaultHeaderMatcher(key)
	}

	return
}

func (s *Server) out(key string) (new string, match bool) {
	if newKey, test := s.config.Gateway.Header.TestOuts(key); test {
		new = newKey
		match = true
	} else {
		new, match = runtime.DefaultHeaderMatcher(key)
	}

	return
}

func (s *Server) metadata(_ context.Context, req *http.Request) metadata.MD {
	md := make(map[string]string)
	md[constant.GrpcGatewayUri] = req.URL.RequestURI()
	md[constant.GrpcGatewayMethod] = req.Method
	md[constant.GrpcGatewayProto] = req.Proto

	return metadata.New(md)
}

func (s *Server) isHttpCode(code int) bool {
	return code >= http.StatusContinue && code <= http.StatusNetworkAuthenticationRequired
}
