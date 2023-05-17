package grpc

import (
	"strings"
)

type raw struct {
	// 前缀
	Prefix string `json:"prefix" yaml:"prefix" xml:"prefix" toml:"prefix"`
	// 后缀
	Suffix string `yaml:"suffix" yaml:"suffix" xml:"suffix" toml:"suffix"`
	// 包含
	Contains string `default:"raw" json:"contains" yaml:"contains" xml:"contains" toml:"contains"`
}

func (r *raw) check(check string) (checked bool) {
	if "" != r.Contains && strings.Contains(check, r.Contains) {
		checked = true
	} else if "" != r.Prefix && strings.HasPrefix(check, r.Prefix) {
		checked = true
	} else if "" != r.Suffix && strings.HasSuffix(check, r.Suffix) {
		checked = true
	}

	return
}
