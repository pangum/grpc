package grpc

type header struct {
	// 删除列表
	Removes []remove `json:"removes" yaml:"removes" xml:"removes" toml:"removes"`
	// 输入头匹配列表
	Ins []matcher `default:"[{'prefix': 'X-Forwarded'}]" json:"ins" yaml:"ins" xml:"ins" toml:"ins"`
	// 输出头匹配列表
	Outs []matcher `xml:"outs" yaml:"outs" xml:"outs" toml:"outs"`
}

func (h *header) testRemove(key string) (new string, match bool) {
	for _, remove := range h.Removes {
		if new, match = remove.test(key); match {
			break
		}
	}

	return
}

func (h *header) testIns(key string) (new string, match bool) {
	return h.match(h.Ins, key)
}

func (h *header) testOuts(key string) (new string, match bool) {
	return h.match(h.Outs, key)
}

func (h *header) match(targets []matcher, key string) (new string, match bool) {
	for _, target := range targets {
		if new, match = target.test(key); match {
			break
		}
	}

	return
}
