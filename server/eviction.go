package server

import "fmt"

type LRU struct{}

func (s *LRU) evict(c *Cache) {
	fmt.Println("Evicting cache using LRU strtegy")
}

func (s *LRU) shouldExpire() bool {
	return false
}

type MRU struct{}

func (s *MRU) evict(c *Cache) {
	fmt.Println("Evicting cache using LRU strtegy")
}

func (s *MRU) shouldExpire() bool {
	return false
}

type LFU struct{}

func (s *LFU) evict(c *Cache) {
	fmt.Println("Evicting cache using LRU strtegy")
}

func (s *LFU) shouldExpire() bool {
	return false
}

type TTL struct{}

func (s *TTL) evict(c *Cache) {
	fmt.Println("Evicting cache using LRU strtegy")
}

func (s *TTL) shouldExpire() bool {
	return true
}
