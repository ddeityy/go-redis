package main

import (
	"context"
	"go-redis/client"
	"go-redis/server"
	"go-redis/server/pb"
	"log"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	go server.RunServer()

	time.Sleep(2 * time.Second)

	c := client.NewClient()
	defer c.Close()

	kv := &pb.KeyValue{Key: "hello", Value: "world"}

	_, err := c.Set(context.Background(), kv)
	if err != nil {
		log.Printf("set: %s", err)
	} else {
		log.Printf("%v cached successfully", kv.Key)
	}

	get, err := c.Get(context.Background(), &pb.Key{Key: "hello"})
	if err != nil {
		log.Printf("get: %s", err)
	} else {
		log.Printf("cached: %s", get.Value)
	}

	wg.Wait()
}
