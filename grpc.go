package main

import (
	"context"
	"google.golang.org/grpc"
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
	)

	if err != nil {
		return nil, err
	}

	return connection, nil
}
