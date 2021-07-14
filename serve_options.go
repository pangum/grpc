package grpc

type serveOptions struct {
	// 启用反射，用于调试
	Reflection bool `default:"true" json:"reflection" yaml:"reflection" xml:"reflection"`
}

func defaultServeOptions() *serveOptions {
	return &serveOptions{
		Reflection: true,
	}
}
