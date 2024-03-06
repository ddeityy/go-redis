package main

import (
	"context"
	"fmt"
	"go-redis/client"
	"go-redis/server"
	"go-redis/server/pb"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func main() {

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	var wg sync.WaitGroup
	wg.Add(20)

	go server.RunServer()

	time.Sleep(2 * time.Second)

	c := client.NewClient()
	defer c.Close()

	go set(100, c, &wg)

	time.Sleep(1 * time.Second)

	go get(10000, c, &wg)
	go get(10000, c, &wg)
	go get(10000, c, &wg)
	go get(10000, c, &wg)
	go get(10000, c, &wg)

	time.Sleep(3 * time.Second)

	go set(100, c, &wg)

	time.Sleep(1 * time.Second)

	go get(10000, c, &wg)
	go get(10000, c, &wg)
	go get(10000, c, &wg)
	go get(10000, c, &wg)
	go get(10000, c, &wg)

	wg.Wait()
}

func status(c *client.Client, wg *sync.WaitGroup) {
	data, err := c.GetCacheData(context.Background(), &pb.Empty{})
	if err != nil {
		log.Println(err)
	}
	for _, v := range data.Items {
		fmt.Println(v.Key)
		fmt.Println(v.Value)
		fmt.Println(v.Accessed)
		fmt.Println("---------------")
	}

	defer wg.Done()
}

func get(n int, c *client.Client, wg *sync.WaitGroup) {
	for range n {
		get, err := c.Get(context.Background(), &pb.Key{Key: fmt.Sprintf("%v", rand.Intn(100))})
		if err != nil {
			log.Printf("get: %s", err)
		} else {
			log.Printf("%s", get.Value)
		}
	}
	defer wg.Done()
}

func set(n int, c *client.Client, wg *sync.WaitGroup) {
	for n := range n {
		kv := &pb.KeyValue{Key: fmt.Sprintf("%v", n), Value: "world"}
		resp, err := c.Set(context.Background(), kv)
		log.Println(resp.Code, resp.Message)
		if err != nil {
			log.Printf("set: %s", err)
		} else {
			log.Printf("%v cached successfully", kv.Key)
		}
	}
	defer wg.Done()
}
