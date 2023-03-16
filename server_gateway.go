package grpc

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/goexl/gox/field"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

func (s *Server) gateway(register register) (err error) {
	if !s.config.gatewayEnabled() {
		return
	}

	gatewayOptions := s.config.Gateway.options()
	gatewayOptions = append(gatewayOptions, runtime.WithForwardResponseOption(s.response))
	gatewayOptions = append(gatewayOptions, runtime.WithIncomingHeaderMatcher(s.in))
	gatewayOptions = append(gatewayOptions, runtime.WithOutgoingHeaderMatcher(s.out))
	gatewayOptions = append(gatewayOptions, runtime.WithMetadata(s.metadata))
	if nil != s.config.Gateway.Unescape {
		gatewayOptions = append(gatewayOptions, runtime.WithUnescapingMode(s.config.Gateway.Unescape.Mode))
	}
	gateway := runtime.NewServeMux(gatewayOptions...)
	grpcOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if ge := register.Gateway(gateway, s.config.Addr(), grpcOptions...); nil != ge {
		err = ge
	} else {
		s.mux.Handle("/", gateway)
		s.serveHttp = true
	}

	return
}

func (s *Server) response(ctx context.Context, writer http.ResponseWriter, msg proto.Message) (err error) {
	// 注意，这儿的顺序不能乱，必须先写入头再写入状态码
	if se := s.header(ctx, writer, msg); nil != se {
		err = se
	} else if he := s.status(ctx, writer, msg); nil != he {
		err = he
	}

	return
}

func (s *Server) status(ctx context.Context, writer http.ResponseWriter, _ proto.Message) (err error) {
	if md, ok := runtime.ServerMetadataFromContext(ctx); !ok {
		// 上下文无法转换
	} else if status := md.HeaderMD.Get(httpStatusHeader); 0 == len(status) {
		// 没有设置状态
	} else if code, ae := strconv.Atoi(status[0]); nil != ae {
		err = ae
		s.logger.Warn("状态码被错误设置", field.New("value", status))
	} else {
		md.HeaderMD.Delete(httpStatusHeader)
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
		if httpStatusHeader == key { // 不处理设置状态码的逻辑，由状态码设置逻辑特殊处理
			continue
		}

		newKey := strings.ToLower(key)
		removal := false
		newKey, removal = s.config.Gateway.Header.testRemove(newKey)

		if removal {
			writer.Header().Set(newKey, strings.Join(value, space))
			header.Delete(key)
			writer.Header().Del(fmt.Sprintf(grpcMetadataFormatter, key))
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

func (s *Server) metadata(_ context.Context, req *http.Request) metadata.MD {
	md := make(map[string]string)
	md[grpcGatewayUri] = req.URL.RequestURI()
	md[grpcGatewayMethod] = req.Method
	md[grpcGatewayProto] = req.Proto

	return metadata.New(md)
}
