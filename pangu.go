package grpc

import (
	"github.com/pangum/grpc/internal/plugin"
	"github.com/pangum/pangu"
)

func init() {
	creator := new(plugin.Constructor)
	pangu.New().Get().Dependency().Put(
		creator.New,
		creator.NewClient,
	).Build().Build().Apply()
}
