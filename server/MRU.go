package server

import (
	"fmt"
)

type MRU struct{}

func (s *MRU) evict(c *Cache) error {
	fmt.Println("Evicting cache using LRU strtegy")
	return nil
}
