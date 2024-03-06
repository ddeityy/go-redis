package server

import (
	"container/list"
	"log"
	"math"
	"time"
)

type LFU struct{}

func (s *LFU) evict(c *Cache) error {
	var leastFrequentItem *list.Element
	var leastFrequency int32 = math.MaxInt32 // Start with the maximum possible frequency
	for e := c.list.Front(); e != nil; e = e.Next() {
		v := e.Value.(value)
		if v.accessed < leastFrequency {
			leastFrequency = v.accessed
			leastFrequentItem = e
		}
	}

	c.store.Delete(leastFrequentItem.Value.(value).key)
	c.list.Remove(leastFrequentItem.Value.(value).elementPtr)
	c.capacity--

	log.Println("lfu evicted", leastFrequentItem.Value.(value).key)

	return nil
}

func (s *LFU) ttl() time.Duration {
	return 0
}
