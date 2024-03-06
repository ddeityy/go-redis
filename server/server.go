package server

import (
	"context"
	"errors"
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
	val, err := s.cache.Get(key.Key)

	if err != nil {
		return nil, err
	}
	log.Println("Get request:", key.Key)

	return &pb.Value{Value: string(val)}, nil
}

func (s Server) Set(ctx context.Context, kv *pb.KeyValue) (*pb.Response, error) {
	log.Println("Set request:", kv.Key)
	if err := s.cache.Set(kv.Key, kv.Value); err != nil {
		if errors.Is(err, fmt.Errorf("set: key already exists")) {
			return &pb.Response{Code: 400, Message: err.Error()}, err
		}
	}

	return &pb.Response{Code: 200, Message: "success"}, nil
}

func (s Server) GetCacheData(ctx context.Context, _ *pb.Empty) (*pb.Data, error) {
	m := s.cache.GetCacheData()
	var data pb.Data
	for k, v := range m {
		data.Items = append(data.Items, &pb.CachedItem{Key: k, Value: v.data, Accessed: int32(v.accessed)})
	}
	return &data, nil
}

func (s Server) SwitchStrategy(ctx context.Context, req *pb.Strategy) (*pb.Response, error) {
	log.Println("SwitchStrategy request:", req.Strategy)

	// Assuming you have a method to set the eviction strategy based on the enum value
	switch req.Strategy {
	case pb.StrategyEnum_LRU:
		s.cache.SetEvictionStrategy(&LRU{})
	case pb.StrategyEnum_LFU:
		s.cache.SetEvictionStrategy(&LFU{})
	case pb.StrategyEnum_TTL:
		s.cache.SetEvictionStrategy(&TTL{})
	default:
		return &pb.Response{Code: 400, Message: "Invalid strategy"}, nil
	}

	return &pb.Response{Code: 200, Message: "Strategy switched successfully"}, nil
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
