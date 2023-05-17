package grpc

type body struct {
	// 原始请求
	Raws []*raw `default:"[{'contains': 'raw'}]" json:"raws" yaml:"raws" xml:"raws" toml:"raws"`
}

func (b *body) check(check string) (checked bool) {
	for _, _raw := range b.Raws {
		if checked = _raw.check(check); checked {
			break
		}
	}

	return
}
