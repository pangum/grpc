package grpc

import (
	"github.com/goexl/gox"
)

type header struct {
	// 是否启用默认行为
	Default *bool `default:"true" json:"default" yaml:"default" xml:"default" toml:"default"`
	// 删除列表
	Removes []remove `json:"removes" yaml:"removes" xml:"removes" toml:"removes"`
	// 输入头匹配列表
	// nolint: lll
	Ins []matcher `json:"ins" yaml:"ins" xml:"ins" toml:"ins"`
	// 输出头匹配列表
	Outs []matcher `xml:"outs" yaml:"outs" xml:"outs" toml:"outs"`

	DefaultRemoves []remove  `default:"[{'prefix': 'http-'}]"`
	DefaultIns     []matcher `default:"[{'prefix': 'x-forwarded'}]"`
}

func (h *header) testRemove(key string) (new string, match bool) {
	for _, remove := range gox.Ifx(*h.Default, append(h.DefaultRemoves, h.Removes...), h.Removes) {
		if new, match = remove.test(key); match {
			break
		}
	}

	return
}

func (h *header) testIns(key string) (new string, match bool) {
	return h.match(gox.Ifx(*h.Default, append(h.DefaultIns, h.Ins...), h.Ins), key)
}

func (h *header) testOuts(key string) (new string, match bool) {
	return h.match(h.Outs, key)
}

func (h *header) match(matchers []matcher, key string) (new string, match bool) {
	for _, matcher := range matchers {
		if new, match = matcher.test(key); match {
			break
		}
	}

	return
}
