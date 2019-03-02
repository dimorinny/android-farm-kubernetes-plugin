package main

import (
	"context"
	"google.golang.org/grpc"
	"net"
	"time"
)

func grpcContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
}

func dial(
	unixSocketPath string,
	ctx context.Context,
) (*grpc.ClientConn, error) {
	connection, err := grpc.DialContext(
		ctx,
		unixSocketPath,
		grpc.WithBlock(),
		grpc.WithInsecure(),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "unix", addr)
		}),
	)
	if err != nil {
		return nil, err
	}

	return connection, nil
}
