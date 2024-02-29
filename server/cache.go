package server

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type value struct {
	Data         string `json:"data"`
	accesed      int
	set          time.Time
	lastAccessed time.Time
	expireAfter  time.Duration
	shouldExpire bool
}

type Cache struct {
	storage      sync.Map
	evictionAlgo EvictionAlgo
	capacity     int
	maxCapacity  int
}

func (c *Cache) Get(key string) (string, error) {
	val, ok := c.storage.Load(key)
	if ok {
		out, ok := val.(value)
		if ok {
			out.accesed++
			out.lastAccessed = time.Now()
			return out.Data, nil
		}
		return "", errors.New("type assertion failed")
	}
	return "", errors.New("key not found")
}

func (c *Cache) Set(key string, val string) {
	c.storage.Store(
		key,
		value{
			Data:         val,
			accesed:      0,
			set:          time.Now(),
			lastAccessed: time.Now(),
			expireAfter:  time.Second * 20,
			shouldExpire: c.evictionAlgo.shouldExpire(),
		},
	)
}

func (c *Cache) Default() *Cache {
	return initCache(&LRU{})
}

func (c *Cache) SetEvictionStrategy(strategy string) error {
	switch strategy {
	case "lru":
		c.evictionAlgo = &LRU{}
	case "mru":
		c.evictionAlgo = &MRU{}
	case "lfu":
		c.evictionAlgo = &LFU{}
	case "ttl":
		c.evictionAlgo = &TTL{}
	default:
		return errors.New("set eviction strategy: strategy not found")
	}

	return nil
}

type EvictionAlgo interface {
	evict(c *Cache)
	shouldExpire() bool
}

func initCache(e EvictionAlgo) *Cache {
	return &Cache{
		storage:      sync.Map{},
		evictionAlgo: e,
		capacity:     0,
		maxCapacity:  2,
	}
}

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
