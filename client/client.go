package client

import (
	"go-redis/server/pb"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	*grpc.ClientConn
	pb.CacheClient
}

func NewClient() *Client {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %s", err)
	}

	c := pb.NewCacheClient(conn)

	return &Client{conn, c}
}
