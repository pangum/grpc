package grpc

import (
	"github.com/goexl/gox"
)

type msg struct {
	// 发送大小
	// 4GB
	Send gox.Size `default:"4GB" json:"send" yaml:"send" xml:"send" toml:"send"`
	// 接收大小
	// 4GB
	Receive gox.Size `default:"4GB" json:"receive" yaml:"receive" xml:"receive" toml:"receive"`
}
