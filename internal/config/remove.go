package config

import (
	"strings"
)

type Remove struct {
	// 前缀
	Prefix string `json:"prefix" yaml:"prefix" xml:"prefix" toml:"suffix"`
	// 后缀
	Suffix string `json:"suffix" yaml:"suffix" xml:"suffix" toml:"suffix"`
}

func (r *Remove) Test(key string) (new string, match bool) {
	key = strings.ToLower(key)
	prefix := r.Prefix
	if "" != prefix && strings.HasPrefix(key, r.Prefix) {
		new = strings.TrimPrefix(key, prefix)
		match = true
	}

	suffix := r.Suffix
	if "" != suffix && strings.HasSuffix(key, suffix) {
		new = strings.TrimSuffix(key, suffix)
		match = true
	}

	return
}
