package client

import (
	"context"
	"go-redis/server"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
}

func RunClient() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := server.NewCacheClient(conn)

	kv := &server.KeyValue{Key: "hello", Value: "world"}

	set, err := c.Set(context.Background(), kv)
	if err != nil {
		log.Fatalf("Error when calling Set: %s", err)
	}
	log.Printf("Cached: %v", set.Stored)

	get, err := c.Get(context.Background(), &server.Key{Key: "hello"})
	if err != nil {
		log.Fatalf("Error when calling Get: %s", err)
	}
	log.Printf("Cached value: %s", get.Value)

}
