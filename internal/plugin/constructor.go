package plugin

import (
	"net/http"

	"github.com/goexl/log"
	"github.com/pangum/grpc/internal/core"
	"github.com/pangum/pangu"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type Constructor struct {
	// 解决命名空间问题
}

func (c *Constructor) New(config *pangu.Config, logger log.Logger) (server *core.Server, mux *http.ServeMux, err error) {
	wrapper := new(Wrapper)
	if ge := config.Build().Get(wrapper); nil != ge {
		err = ge
	} else {
		server, mux, err = c.new(wrapper.Grpc, logger)
	}

	return
}

func (c *Constructor) NewClient(config *pangu.Config) (client *core.Client, err error) {
	wrapper := new(Wrapper)
	if ge := config.Build().Get(wrapper); nil != ge {
		err = ge
	} else {
		client, err = c.newClient(wrapper.Grpc)
	}

	return
}

func (c *Constructor) NewGateway(logger log.Logger) *core.Gateway {
	return core.NewGateway(logger)
}

func (c *Constructor) new(config *Config, logger log.Logger) (server *core.Server, mux *http.ServeMux, err error) {
	if nil == config.Server {
		return
	}

	opts := make([]grpc.ServerOption, 0, 8)
	opts = append(opts, grpc.InitialWindowSize(int32(config.Options.Size.Window.Initial)))
	opts = append(opts, grpc.InitialConnWindowSize(int32(config.Options.Size.Window.Connection)))
	opts = append(opts, grpc.MaxSendMsgSize(int(config.Options.Size.Msg.Send)))
	opts = append(opts, grpc.MaxRecvMsgSize(int(config.Options.Size.Msg.Receive)))
	opts = append(opts, grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		PermitWithoutStream: config.Options.Keepalive.Policy.Permit,
	}))
	opts = append(opts, grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: config.Options.Keepalive.Idle,
		Time:              config.Options.Keepalive.Time,
		Timeout:           config.Options.Keepalive.Timeout,
	}))

	mux = http.NewServeMux()
	server = core.NewServer(grpc.NewServer(opts...), mux, config.Server, config.Gateway, logger)

	return
}

func (c *Constructor) newClient(config *Config) (client *core.Client, err error) {
	if 0 == len(config.Clients) {
		return
	}

	opts := make([]grpc.DialOption, 0, 8)
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithInitialWindowSize(int32(config.Options.Size.Window.Initial)))
	opts = append(opts, grpc.WithInitialConnWindowSize(int32(config.Options.Size.Window.Connection)))
	opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(int(config.Options.Size.Msg.Send))))
	opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(int(config.Options.Size.Msg.Receive))))
	opts = append(opts, grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                config.Options.Keepalive.Time,
		Timeout:             config.Options.Keepalive.Timeout,
		PermitWithoutStream: config.Options.Keepalive.Policy.Permit,
	}))

	connections := make(map[string]*grpc.ClientConn)
	for _, conf := range config.Clients {
		var connection *grpc.ClientConn
		if connection, err = grpc.Dial(conf.Addr, opts...); nil != err {
			return
		}

		if "" != conf.Name {
			connections[conf.Name] = connection
		}
		for _, name := range conf.Names {
			connections[name] = connection
		}
	}
	client = core.NewClient(connections)

	return
}
