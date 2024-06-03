package config

import (
	"github.com/goexl/gox"
)

type Msg struct {
	// 发送大小
	// 4GB
	Send gox.Bytes `default:"4GB" json:"send" yaml:"send" xml:"send" toml:"send"`
	// 接收大小
	// 4GB
	Receive gox.Bytes `default:"4GB" json:"receive" yaml:"receive" xml:"receive" toml:"receive"`
}
