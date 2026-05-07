package clients

import (
	"mangahub/proto/chapter"
	"mangahub/proto/manga"
	"mangahub/proto/message"
	"mangahub/proto/session"
	"mangahub/proto/user"
	"mangahub/proto/user_manga"

	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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

func NewMangaGRPCClient() (manga.GRPCMangaServiceClient, *grpc.ClientConn, error) {
    // Define gRPC server address
	target := os.Getenv("SERVER_HOST") + ":" + os.Getenv("GRPC_SERVER_PORT")
    conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return nil, nil, err
    }

    // return client & connection
    client := manga.NewGRPCMangaServiceClient(conn)
    return client, conn, nil
}

func NewUserMangaGRPCClient() (user_manga.GRPCUserMangaServiceClient, *grpc.ClientConn, error) {
    // Define gRPC server address
	target := os.Getenv("SERVER_HOST") + ":" + os.Getenv("GRPC_SERVER_PORT")
    conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return nil, nil, err
    }

    // return client & connection
    client := user_manga.NewGRPCUserMangaServiceClient(conn)
    return client, conn, nil
}

func NewChapterGRPCClient() (chapter.GRPCChapterServiceClient, *grpc.ClientConn, error) {
	// Define gRPC server address
	target := os.Getenv("SERVER_HOST") + ":" + os.Getenv("GRPC_SERVER_PORT")
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	// return client & connection
	client := chapter.NewGRPCChapterServiceClient(conn)
	return client, conn, nil
}

func NewMessageGRPCClient() (message.GRPCMessageServiceClient, *grpc.ClientConn, error) {
	// Define gRPC server address
	target := os.Getenv("SERVER_HOST") + ":" + os.Getenv("GRPC_SERVER_PORT")
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	// return client & connection
	client := message.NewGRPCMessageServiceClient(conn)
	return client, conn, nil
}
