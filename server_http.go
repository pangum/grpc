package grpc

import (
	"net"
	"net/http"
	"strings"

	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

func (s *Server) handler(grpc *grpc.Server, gateway http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		s.Debug("收到请求", s.fields(request)...)
		if request.ProtoMajor >= 2 && grpcHeaderValue == request.Header.Get(headerContentType) {
			grpc.ServeHTTP(writer, request)
		} else {
			s.addRawType(request)
			gateway.ServeHTTP(writer, request)
		}
	}), new(http2.Server))
}

func (s *Server) addRawType(request *http.Request) {
	if nil != request && s.config.Gateway.Body.check(request.URL.Path) {
		request.Header.Set(headerContentType, rawHeaderValue)
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
	ip = req.Header.Get(headerXRealIp)
	if netIp := net.ParseIP(ip); nil != netIp {
		return
	}

	ips := req.Header.Get(headerXForwardedFor)
	for _, _ip := range strings.Split(ips, comma) {
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
