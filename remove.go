package grpc

import (
	"strings"
)

type remove struct {
	// 前缀
	Prefix string `default:"Http-" json:"prefix" yaml:"prefix" xml:"prefix" toml:"suffix"`
	// 后缀
	Suffix string `json:"suffix" yaml:"suffix" xml:"suffix" toml:"suffix"`
}

func (r *remove) test(key string) (new string, match bool) {
	key = strings.ToLower(key)
	prefix := strings.ToLower(r.Prefix)
	if strings.HasPrefix(key, prefix) {
		new = strings.TrimPrefix(key, prefix)
		match = true
	}

	suffix := strings.ToLower(r.Suffix)
	if "" != suffix && strings.HasSuffix(key, suffix) {
		new = strings.TrimSuffix(key, suffix)
		match = true
	}

	return
}
