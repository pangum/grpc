package grpc

import (
	"github.com/goexl/gox"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

type gateway struct {
	// 是否开启
	Enabled *bool `default:"true" json:"enabled" yaml:"enabled" xml:"enabled" toml:"enabled"`
	// 路径
	Path string `json:"path" yaml:"path" xml:"path" toml:"path" validate:"omitempty,startswith=/,endsnotwith=/"`
	// 序列化
	Json json `json:"json" yaml:"json" xml:"json" toml:"json"`
	// 头
	Header header `json:"header" yaml:"header" xml:"header" toml:"header"`
	// 消息体
	Body body `json:"body" yaml:"body" xml:"body" toml:"body"`
	// 模式
	Unescape *unescape `json:"unescape" yaml:"unescape" xml:"unescape" toml:"unescape"`
}

func (g *gateway) options() (options []runtime.ServeMuxOption) {
	options = make([]runtime.ServeMuxOption, 0, 1)
	options = append(options, runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			Multiline:       g.Json.Multiline,
			Indent:          g.Json.Indent,
			AllowPartial:    g.Json.Partial,
			UseProtoNames:   gox.Contains(&g.Json.Options, nameAsProto),
			UseEnumNumbers:  gox.Contains(&g.Json.Options, enumAsNumbers),
			EmitUnpopulated: *g.Json.Unpopulated,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			AllowPartial:   g.Json.Partial,
			DiscardUnknown: *g.Json.Discard,
		},
	}))

	return
}
