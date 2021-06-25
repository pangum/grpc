package grpc

import `github.com/storezhang/pangu`

func init() {
	if err := pangu.New().Provides(New); nil != err {
		panic(err)
	}
}
