package grpc

import (
	"github.com/goexl/gox"
	"github.com/pangum/grpc/internal/core"
	"github.com/pangum/pangu"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// Client gRPC客户端封装
type Client struct {
	connections map[string]*grpc.ClientConn

	_ gox.CannotCopy
}

func newClient(config *pangu.Config) (client *Client, err error) {
	wrap := new(core.Wrapper)
	if err = config.Build().Get(wrap); nil != err {
		return
	}

	// 组织配置项
	conf := wrap.Grpc
	opts := make([]grpc.DialOption, 0, 8)
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithInitialWindowSize(int32(conf.Options.Size.Window.Initial)))
	opts = append(opts, grpc.WithInitialConnWindowSize(int32(conf.Options.Size.Window.Connection)))
	opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(int(conf.Options.Size.Msg.Send))))
	opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(int(conf.Options.Size.Msg.Receive))))
	opts = append(opts, grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                conf.Options.Keepalive.Time,
		Timeout:             conf.Options.Keepalive.Timeout,
		PermitWithoutStream: conf.Options.Keepalive.Policy.Permit,
	}))

	connections := make(map[string]*grpc.ClientConn)
	for _, conf := range conf.Clients {
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
	client = &Client{
		connections: connections,
	}

	return
}

func (c *Client) Connection(name string) *grpc.ClientConn {
	return c.connections[name]
}
