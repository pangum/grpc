package core

import (
	"github.com/goexl/gox"
	"google.golang.org/grpc"
)

// Client gRPC客户端封装
type Client struct {
	connections map[string]*grpc.ClientConn

	_ gox.CannotCopy
}

func NewClient(connections map[string]*grpc.ClientConn) *Client {
	return &Client{
		connections: connections,
	}
}

func (c *Client) Connection(name string) *grpc.ClientConn {
	return c.connections[name]
}
