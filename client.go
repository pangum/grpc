package grpc

import (
	"github.com/goexl/gox"
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
	_panguConfig := new(panguConfig)
	if err = config.Load(_panguConfig); nil != err {
		return
	}

	// 组织配置项
	_config := _panguConfig.Grpc
	_options := make([]grpc.DialOption, 0, 8)
	_options = append(_options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	_options = append(_options, grpc.WithInitialWindowSize(int32(_config.Options.Size.Window.Initial)))
	_options = append(_options, grpc.WithInitialConnWindowSize(int32(_config.Options.Size.Window.Connection)))
	// nolint:lll
	_options = append(_options, grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(int(_config.Options.Size.Msg.Send))))
	// nolint:lll
	_options = append(_options, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(int(_config.Options.Size.Msg.Receive))))
	_options = append(_options, grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                _config.Options.Keepalive.Time,
		Timeout:             _config.Options.Keepalive.Timeout,
		PermitWithoutStream: _config.Options.Keepalive.Policy.Permit,
	}))

	connections := make(map[string]*grpc.ClientConn)
	for _, conf := range _config.Clients {
		var connection *grpc.ClientConn
		if connection, err = grpc.Dial(conf.Addr, _options...); nil != err {
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
