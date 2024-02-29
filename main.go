package main

import (
	"go-redis/client"
	"go-redis/server"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go server.RunServer()

	time.Sleep(1 * time.Second)

	go client.RunClient()

	wg.Wait()
}
