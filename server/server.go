package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"go-redis/server/pb"

	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedCacheServer
	cache *Cache
}

func (s Server) Get(ctx context.Context, key *pb.Key) (*pb.Value, error) {
	log.Println("Get request:", key.Key)
	val, err := s.cache.Get(key.Key)

	if err != nil {
		return nil, err
	}

	return &pb.Value{Value: string(val)}, nil
}

func (s Server) Set(ctx context.Context, kv *pb.KeyValue) (*pb.Cached, error) {
	log.Println("Set request:", kv.Key)
	set := s.cache.Set(kv.Key, kv.Value)
	if set {
		return &pb.Cached{Cached: true}, nil
	}
	return nil, fmt.Errorf("key %v already exists", kv.Key)
}

func RunServer() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	cache := NewCache()
	serv := Server{pb.UnimplementedCacheServer{}, cache}
	pb.RegisterCacheServer(grpcServer, serv)

	log.Println("Listening on", lis.Addr().String())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
