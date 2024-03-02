package server

import (
	"container/list"
	"fmt"
)

type LRU struct {
	list *list.List
}

func (s *LRU) evict(c *Cache) error {
	fmt.Println("Evicting cache using LRU strtegy")

	last := s.list.Back()
	if last == nil {
		return fmt.Errorf("lru evict: can not evict from a nil list")
	}

	s.list.Remove(last)

	c.storage.Delete(last.Value)

	return nil
}
