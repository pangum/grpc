package grpc

import (
	"strings"
)

type matcher struct {
	// 等于
	Equal string `json:"equal" yaml:"equal" xml:"equal" toml:"equal"`
	// 前缀
	Prefix string `json:"prefix" yaml:"prefix" xml:"prefix" toml:"suffix"`
	// 后缀
	Suffix string `json:"suffix" yaml:"suffix" xml:"suffix" toml:"suffix"`
	// 包含
	Contains string `json:"contains" yaml:"contains" xml:"contains" toml:"contains"`
}

func (m *matcher) test(key string) (new string, match bool) {
	key = strings.ToLower(key)
	new = key
	if "" != m.Equal && m.Equal == key {
		match = true
	}

	if "" != m.Prefix && strings.HasPrefix(key, m.Prefix) {
		match = true
	}

	if "" != m.Suffix && strings.HasSuffix(key, m.Suffix) {
		new = key
		match = true
	}

	if "" != m.Contains && strings.Contains(key, m.Contains) {
		new = key
		match = true
	}

	return
}
