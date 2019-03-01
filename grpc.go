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

// TODO: use dial context instead deprecated dial
func dial(
	unixSocketPath string,
	ctx context.Context,
) (*grpc.ClientConn, error) {
	//connection, err := grpc.DialContext(
	//	ctx,
	//	unixSocketPath,
	//	grpc.WithBlock(),
	//	grpc.WithInsecure(),
	//)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return connection, nil

	connection, err := grpc.Dial(
		unixSocketPath,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second*5),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	)

	if err != nil {
		return nil, err
	}

	return connection, nil
}
