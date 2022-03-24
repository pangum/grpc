package grpc

import (
	`github.com/goexl/gox`
	`github.com/pangum/pangu`
	`google.golang.org/grpc`
	`google.golang.org/grpc/keepalive`
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
	grpcConfig := _panguConfig.Grpc
	_options := make([]grpc.DialOption, 0, 8)
	_options = append(_options, grpc.WithInsecure())
	_options = append(_options, grpc.WithInitialWindowSize(grpcConfig.Options.Size.Window.Initial))
	_options = append(_options, grpc.WithInitialConnWindowSize(grpcConfig.Options.Size.Window.Connection))
	_options = append(_options, grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(grpcConfig.Options.Size.Msg.Send)))
	_options = append(_options, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcConfig.Options.Size.Msg.Receive)))
	_options = append(_options, grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                grpcConfig.Options.Keepalive.Time,
		Timeout:             grpcConfig.Options.Keepalive.Timeout,
		PermitWithoutStream: grpcConfig.Options.Keepalive.Policy.Permit,
	}))

	connections := make(map[string]*grpc.ClientConn)
	for _, _clientConfig := range grpcConfig.Clients {
		var connection *grpc.ClientConn
		if connection, err = grpc.Dial(_clientConfig.Addr, _options...); nil != err {
			return
		}

		if "" != _clientConfig.Name {
			connections[_clientConfig.Name] = connection
		}
		for _, name := range _clientConfig.Names {
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
