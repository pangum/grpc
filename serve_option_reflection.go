package grpc

var (
	_             = DisableReflection
	_ serveOption = (*serveOptionReflection)(nil)
)

type serveOptionReflection struct{}

// DisableReflection 禁用反射
func DisableReflection() *serveOptionReflection {
	return &serveOptionReflection{}
}

func (r *serveOptionReflection) apply(options *serveOptions) {
	options.Reflection = false
}
