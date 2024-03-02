package server

import (
	"fmt"
)

type TTL struct{}

func (s *TTL) evict(c *Cache) error {
	fmt.Println("Evicting cache using LRU strtegy")
	return nil
}
