package config

import (
	"fmt"

	"github.com/goexl/gox"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pangum/grpc/internal/internal/constant"
	"google.golang.org/protobuf/encoding/protojson"
)

type Gateway struct {
	// 是否开启
	Enabled *bool `default:"true" json:"enabled" yaml:"enabled" xml:"enabled" toml:"enabled"`
	// 绑定监听主机
	Host string `json:"host" yaml:"host" xml:"host" toml:"host"`
	// 绑定监听端口
	Port int `default:"9001" json:"port" yaml:"port" xml:"port" toml:"port" validate:"required,min=1,max=65535"`
	// 路径
	Path string `json:"path" yaml:"path" xml:"path" toml:"path" validate:"omitempty,startswith=/,endsnotwith=/"`
	// 跨域
	Cors *Cors `json:"cors" yaml:"cors" xml:"cors" toml:"cors"`
	// 超时
	Timeout Timeout `json:"timeout" yaml:"timeout" xml:"timeout" toml:"timeout"`
	// 序列化
	Json Json `json:"json" yaml:"json" xml:"json" toml:"json"`
	// 头
	Header Header `json:"header" yaml:"header" xml:"header" toml:"header"`
	// 消息体
	Body Body `json:"body" yaml:"body" xml:"body" toml:"body"`
	// 模式
	Unescape *Unescape `json:"unescape" yaml:"unescape" xml:"unescape" toml:"unescape"`
}

func (g *Gateway) Options() (options []runtime.ServeMuxOption) {
	options = make([]runtime.ServeMuxOption, 0, 1)
	json := &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			Multiline:       g.Json.Multiline,
			Indent:          g.Json.Indent,
			AllowPartial:    g.Json.Partial,
			UseProtoNames:   gox.Contains(&g.Json.Options, constant.NameAsProto),
			UseEnumNumbers:  gox.Contains(&g.Json.Options, constant.EnumAsNumbers),
			EmitUnpopulated: g.Json.Unpopulated,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			AllowPartial:   g.Json.Partial,
			DiscardUnknown: *g.Json.Discard,
		},
	}
	options = append(options, runtime.WithMarshalerOption(runtime.MIMEWildcard, json))

	return
}

func (g *Gateway) Addr() string {
	return fmt.Sprintf("%s:%d", g.Host, g.Port)
}

func (g *Gateway) Enable() bool {
	return nil != g.Enabled && *g.Enabled
}

func (g *Gateway) CorsEnabled() bool {
	return nil != g.Cors && *g.Cors.Enabled
}
