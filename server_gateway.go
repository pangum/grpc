package grpc

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/goexl/gox/field"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

func (s *Server) setupGateway(register register) (err error) {
	if !s.config.gatewayEnabled() {
		return
	}

	gatewayOptions := s.config.Gateway.options()
	gatewayOptions = append(gatewayOptions, runtime.WithForwardResponseOption(s.response))
	gatewayOptions = append(gatewayOptions, runtime.WithIncomingHeaderMatcher(s.in))
	gatewayOptions = append(gatewayOptions, runtime.WithOutgoingHeaderMatcher(s.out))
	if nil != s.config.Gateway.Unescape {
		gatewayOptions = append(gatewayOptions, runtime.WithUnescapingMode(s.config.Gateway.Unescape.Mode))
	}
	gateway := runtime.NewServeMux(gatewayOptions...)
	grpcOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err = register.Gateway(gateway, s.config.Addr(), grpcOptions...); nil != err {
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

func (s *Server) response(ctx context.Context, writer http.ResponseWriter, msg proto.Message) (err error) {
	if se := s.status(ctx, writer, msg); nil != se {
		err = se
	} else if he := s.header(ctx, writer, msg); nil != he {
		err = he
	}

	return
}

func (s *Server) status(ctx context.Context, writer http.ResponseWriter, _ proto.Message) (err error) {
	var header metadata.MD
	if md, ok := runtime.ServerMetadataFromContext(ctx); ok {
		header = md.HeaderMD
	}

	status := header.Get(httpStatusHeader)
	if 0 == len(status) {
		s.logger.Warn("状态码被错误设置", field.New("value", status))
	} else if code, ae := strconv.Atoi(status[0]); nil != ae {
		err = ae
		s.logger.Warn("状态码被错误设置", field.New("value", status[0]))
	} else {
		header.Delete(httpStatusHeader)
		writer.Header().Del(grpcStatusHeader)
		writer.WriteHeader(code)
	}

	return
}

func (s *Server) header(ctx context.Context, writer http.ResponseWriter, _ proto.Message) (err error) {
	var header metadata.MD
	if md, ok := runtime.ServerMetadataFromContext(ctx); ok {
		header = md.HeaderMD
	}

	for key, value := range header {
		newKey := strings.ToLower(key)
		removal := false
		newKey, removal = s.config.Gateway.Header.testRemove(newKey)

		if removal {
			writer.Header().Set(newKey, value[0])
			header.Delete(key)
			writer.Header().Del(fmt.Sprintf(grpcMeatadataFormatter, key))
		}
	}

	return
}

func (s *Server) in(key string) (new string, match bool) {
	if newKey, test := s.config.Gateway.Header.testIns(key); test {
		new = newKey
		match = true
	} else {
		new, match = runtime.DefaultHeaderMatcher(key)
	}

	return
}

func (s *Server) out(key string) (new string, match bool) {
	if newKey, test := s.config.Gateway.Header.testOuts(key); test {
		new = newKey
		match = true
	} else {
		new, match = runtime.DefaultHeaderMatcher(key)
	}

	return
}
