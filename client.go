package grpc

import (
	`github.com/storezhang/pangu`
	`google.golang.org/grpc`
	`google.golang.org/grpc/keepalive`
)

// Client gRPC客户端封装
type Client struct {
	connections map[string]*grpc.ClientConn
}

func newClient(config *pangu.Config) (client *Client, err error) {
	panguConfig := new(panguConfig)
	if err = config.Load(panguConfig); nil != err {
		return
	}

	// 组织配置项
	grpcConfig := panguConfig.Grpc
	options := make([]grpc.DialOption, 0, 8)
	options = append(options, grpc.WithInsecure())
	options = append(options, grpc.WithInitialWindowSize(grpcConfig.Options.Size.Window.Initial))
	options = append(options, grpc.WithInitialConnWindowSize(grpcConfig.Options.Size.Window.Connection))
	options = append(options, grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(grpcConfig.Options.Size.Msg.Send)))
	options = append(options, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcConfig.Options.Size.Msg.Receive)))
	options = append(options, grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                grpcConfig.Options.Keepalive.Time,
		Timeout:             grpcConfig.Options.Keepalive.Timeout,
		PermitWithoutStream: grpcConfig.Options.Keepalive.Policy.Permit,
	}))

	connections := make(map[string]*grpc.ClientConn)
	for _, clientConfig := range grpcConfig.Clients {
		var connection *grpc.ClientConn
		if connection, err = grpc.Dial(clientConfig.Addr, options...); nil != err {
			return
		}
		connections[clientConfig.Name] = connection
	}
	client = &Client{
		connections: connections,
	}

	return
}

func (c *Client) Connection(name string) *grpc.ClientConn {
	return c.connections[name]
}
