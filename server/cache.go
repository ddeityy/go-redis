package server

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type value struct {
	data         string
	accesed      int
	set          time.Time
	lastAccessed time.Time
	expireAfter  time.Duration
}

type EvictionAlgo interface {
	evict(*Cache) error
}

type Eviction struct {
	EvictionAlgo
}

type Cache struct {
	storage      sync.Map
	evictionAlgo EvictionAlgo
	capacity     int
	maxCapacity  int
}

// Returns a new cache instance with LRU by default
func NewCache() *Cache {
	return &Cache{
		storage:      sync.Map{},
		evictionAlgo: &LRU{},
		capacity:     0,
		maxCapacity:  500,
	}
}

func (c *Cache) Get(key string) (string, error) {
	val, ok := c.storage.Load(key)
	if ok {
		out, ok := val.(value)
		if ok {
			out.accesed++
			out.lastAccessed = time.Now()
			return out.data, nil
		}
		return "", errors.New("type assertion failed")
	}
	return "", errors.New("key not found")
}

// Puts a key-value pair into cache, returns false if key already exists
func (c *Cache) Set(key string, val string) error {

	if c.capacity >= c.maxCapacity {
		if err := c.evictionAlgo.evict(c); err != nil {
			return err
		}
	}

	if _, ok := c.storage.Load(key); ok {
		return fmt.Errorf("set: key already exists")
	}

	c.storage.Store(
		key,
		value{
			data:         val,
			accesed:      0,
			set:          time.Now(),
			lastAccessed: time.Now(),
			expireAfter:  time.Second * 20,
		},
	)

	return nil
}

func (c *Cache) SetEvictionStrategy(algo EvictionAlgo) {
	c.evictionAlgo = algo
}
