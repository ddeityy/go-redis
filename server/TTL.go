package server

import "time"

type TTL struct{}

func (s *TTL) evict(c *Cache) error {
	return nil
}

func (s *TTL) ttl() time.Duration {
	ttl := time.Second * 20
	return ttl
}
