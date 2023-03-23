package grpc

import (
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

func (s *Server) handler(grpc *grpc.Server, gateway http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor >= 2 && grpcHeaderValue == r.Header.Get(headerContentType) {
			grpc.ServeHTTP(w, r)
		} else {
			gateway.ServeHTTP(w, r)
		}
	}), new(http2.Server))
}
