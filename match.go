package grpc

import (
	"strings"
)

type match struct {
	// 前缀
	Prefix string `json:"prefix" yaml:"prefix" xml:"prefix" toml:"suffix"`
	// 后缀
	Suffix string `json:"suffix" yaml:"suffix" xml:"suffix" toml:"suffix"`
}

func (m *match) test(key string) (new string, match bool) {
	key = strings.ToLower(key)
	prefix := strings.ToLower(m.Prefix)
	if strings.HasPrefix(key, prefix) {
		new = key
		match = true
	}

	suffix := strings.ToLower(m.Suffix)
	if "" != suffix && strings.HasSuffix(key, suffix) {
		new = key
		match = true
	}

	return
}
