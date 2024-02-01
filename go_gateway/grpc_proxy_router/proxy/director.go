package proxy

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type StreamDirector func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error)
