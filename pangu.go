package grpc

import (
	"github.com/pangum/grpc/internal/plugin"
	"github.com/pangum/pangu"
)

func init() {
	ctor := new(plugin.Constructor)
	pangu.New().Get().Dependency().Put(
		ctor.New,
		ctor.NewClient,
		ctor.NewGateway,
		ctor.NewException,
	).Build().Build().Apply()
}
