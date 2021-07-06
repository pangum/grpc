package grpc

import `github.com/storezhang/pangu`

func init() {
	if err := pangu.New().Provides(newServer, newClient); nil != err {
		panic(err)
	}
}
