package server

import (
	"fmt"
)

type LFU struct{}

func (s *LFU) evict(c *Cache) error {
	fmt.Println("Evicting cache using LRU strtegy")
	return nil
}
