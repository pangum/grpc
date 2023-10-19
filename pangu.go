package grpc

import "github.com/pangum/pangu"

func init() {
	pangu.New().Get().Dependencies().Build().Provide(
		newServer,
		newClient,
	)
}
