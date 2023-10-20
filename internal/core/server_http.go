package core

import (
	"net"
	"net/http"
	"strings"

	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/pangum/grpc/internal"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

func (s *Server) handler(grpc *grpc.Server, gateway http.Handler) (handler http.Handler) {
	combine := s.combine(grpc, gateway)
	handler = gox.Ift(s.config.CorsEnabled(), s.cors(combine), combine)

	return
}

func (s *Server) combine(grpc *grpc.Server, gateway http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(rsp http.ResponseWriter, req *http.Request) {
		s.logger.Debug("收到请求", s.fields(req)...)
		if req.ProtoMajor >= 2 && internal.GrpcHeaderValue == req.Header.Get(internal.HeaderContentType) {
			grpc.ServeHTTP(rsp, req)
		} else {
			s.addRawType(req)
			gateway.ServeHTTP(rsp, req)
		}
	}), new(http2.Server))
}

func (s *Server) cors(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rsp http.ResponseWriter, req *http.Request) {
		// 设置允许跨域访问的源
		rsp.Header().Set(internal.HeaderAllowOrigin, strings.Join(s.config.Cors.Allows, internal.Comma))
		// 设置允许的请求方法
		rsp.Header().Set(internal.HeaderAllowMethods, strings.Join(s.config.Cors.Methods, internal.Comma))
		// 设置允许的请求头
		rsp.Header().Set(internal.HeaderAllowHeaders, strings.Join(s.config.Cors.Headers, internal.Comma))

		// 如果是预检请求，直接返回
		if req.Method == internal.MethodOptions {
			return
		}

		// 调用实际的处理器函数
		handler.ServeHTTP(rsp, req)
	})
}

func (s *Server) addRawType(request *http.Request) {
	if nil != request && nil != s.config.Gateway && s.config.Gateway.Body.Check(request.URL.Path) {
		request.Header.Set(internal.HeaderContentType, internal.RawHeaderValue)
	}
}

func (s *Server) fields(request *http.Request) gox.Fields[any] {
	return gox.Fields[any]{
		field.New("method", request.Method),
		field.New("url", request.URL.String()),
		field.New("ip", s.ip(request)),
		field.New("useragent", request.UserAgent()),
		field.New("referer", request.Referer()),
	}
}

func (s *Server) ip(req *http.Request) (ip string) {
	ip = req.Header.Get(internal.HeaderXRealIp)
	if netIp := net.ParseIP(ip); nil != netIp {
		return
	}

	ips := req.Header.Get(internal.HeaderXForwardedFor)
	for _, _ip := range strings.Split(ips, internal.Comma) {
		ip = _ip
		if netIP := net.ParseIP(ip); nil != netIP {
			return
		}
	}

	if _ip, _, err := net.SplitHostPort(req.RemoteAddr); nil == err {
		if netIp := net.ParseIP(_ip); nil != netIp {
			ip = _ip
		}
	}

	return
}
