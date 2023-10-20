package grpc

import (
	"github.com/pangum/grpc/internal/plugin"
	"github.com/pangum/pangu"
)

func init() {
	creator := new(plugin.Creator)
	pangu.New().Get().Dependencies().Build().Provide(
		creator.New,
		creator.NewClient,
	)
}
