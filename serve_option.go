package grpc

type serveOption interface {
	apply(options *serveOptions)
}
