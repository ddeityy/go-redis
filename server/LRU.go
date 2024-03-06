package server

import (
	"fmt"
	"log"
	"time"
)

type LRU struct{}

func (s *LRU) evict(c *Cache) error {
	if c.list.Len() == 0 {
		return fmt.Errorf("empty list")
	}
	var key string
	val, ok := c.list.Front().Value.(value)
	if ok {
		key = val.key
	} else {
		return fmt.Errorf("nil element")
	}
	c.store.Delete(key)
	c.list.Remove(c.list.Front())
	c.capacity--
	log.Println("lru evicted", key)
	return nil
}

func (s *LRU) ttl() time.Duration {
	return 0
}
