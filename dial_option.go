package grpc

import (
	"google.golang.org/grpc"
)

// DialOption 纯粹是为了解决包冲突
type DialOption = grpc.DialOption
