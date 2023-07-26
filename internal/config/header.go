package config

import (
	"github.com/goexl/gox"
)

type Header struct {
	// 是否启用默认行为
	Default *bool `default:"true" json:"default" yaml:"default" xml:"default" toml:"default"`
	// 删除列表
	Removes []Remove `json:"removes" yaml:"removes" xml:"removes" toml:"removes"`
	// 输入头匹配列表
	// nolint: lll
	Ins []Matcher `json:"ins" yaml:"ins" xml:"ins" toml:"ins"`
	// 输出头匹配列表
	Outs []Matcher `json:"outs" yaml:"outs" xml:"outs" toml:"outs"`

	DefaultRemoves []Remove  `default:"[{'prefix': 'http-'}]"`
	DefaultIns     []Matcher `default:"[{'prefix': 'x-forwarded'}]"`
}

func (h *Header) TestRemove(key string) (new string, match bool) {
	removes := gox.Ifx(*h.Default, func() []Remove {
		return append(h.DefaultRemoves, h.Removes...)
	}, func() []Remove {
		return h.Removes
	})
	for _, remove := range removes {
		if new, match = remove.Test(key); match {
			break
		}
	}

	return
}

func (h *Header) TestIns(key string) (new string, match bool) {
	return h.match(gox.Ifx(*h.Default, func() []Matcher {
		return append(h.DefaultIns, h.Ins...)
	}, func() []Matcher {
		return h.Ins
	}), key)
}

func (h *Header) TestOuts(key string) (new string, match bool) {
	return h.match(h.Outs, key)
}

func (h *Header) match(matchers []Matcher, key string) (new string, match bool) {
	for _, matcher := range matchers {
		if new, match = matcher.Test(key); match {
			break
		}
	}

	return
}
