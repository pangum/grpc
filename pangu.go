package grpc

import `github.com/pangum/pangu`

func init() {
	if err := pangu.New().Provides(newServer, newClient); nil != err {
		panic(err)
	}
}
