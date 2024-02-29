package server

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	cache *Cache
}

// mustEmbedUnimplementedCacheServer implements CacheServer.
func (s Server) mustEmbedUnimplementedCacheServer() {
	panic("unimplemented")
}

func (s Server) Get(ctx context.Context, key *Key) (*Value, error) {
	log.Println("Get request:", key.Key)
	val, err := s.cache.Get(key.Key)

	if err != nil {
		return nil, err
	}

	return &Value{Value: string(val)}, nil
}

func (s Server) Set(ctx context.Context, kv *KeyValue) (*Stored, error) {
	log.Println("Set request:", kv.Key)
	s.cache.Set(kv.Key, kv.Value)
	return &Stored{Stored: true}, nil
}

func RunServer() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	cache := Cache{}
	serv := Server{cache.Default()}
	RegisterCacheServer(grpcServer, serv)

	log.Println("Listening on", lis.Addr().String())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
