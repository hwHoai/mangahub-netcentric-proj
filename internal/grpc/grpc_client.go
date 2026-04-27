package grpc

import (
	"mangahub/proto/user"
	"mangahub/proto/session"

	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// internal/client/grpc_client.go
func NewUserGRPCClient() (user.GRPCUserServiceClient, *grpc.ClientConn, error) {
    // Define gRPC server address
	target := os.Getenv("SERVER_HOST") + ":" + os.Getenv("GRPC_SERVER_PORT")
    conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return nil, nil, err
    }

    // return client & connection
    client := user.NewGRPCUserServiceClient(conn)
    return client, conn, nil
}
func NewSessionGRPCClient() (session.GRPCSessionServiceClient, *grpc.ClientConn, error) {
    // Define gRPC server address
	target := os.Getenv("SERVER_HOST") + ":" + os.Getenv("GRPC_SERVER_PORT")
    conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return nil, nil, err
    }

    // return client & connection
    client := session.NewGRPCSessionServiceClient(conn)
    return client, conn, nil
}
