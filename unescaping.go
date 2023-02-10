package grpc

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type unescape struct {
	// 模式
	Mode runtime.UnescapingMode `json:"mode" yaml:"mode" xml:"mode" toml:"mode" validate:"max=3"`
}
