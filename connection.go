package grpc

import (
	"google.golang.org/grpc"
)

// Connection 解决包冲突
type Connection = grpc.ClientConn
