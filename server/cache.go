package server

import (
	"errors"
	"sync"
	"time"
)

type value struct {
	data         string
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
			return out.data, nil
		}
		return "", errors.New("type assertion failed")
	}
	return "", errors.New("key not found")
}

// Puts a key-value pair into cache, returns false if key already exists
func (c *Cache) Set(key string, val string) bool {
	_, loaded := c.storage.LoadOrStore(
		key,
		value{
			data:         val,
			accesed:      0,
			set:          time.Now(),
			lastAccessed: time.Now(),
			expireAfter:  time.Second * 20,
			shouldExpire: c.evictionAlgo.shouldExpire(),
		},
	)
	return !loaded
}

func NewCache() *Cache {
	return initCache(&LRU{})
}

type EvictionAlgo interface {
	evict(c *Cache)
	shouldExpire() bool
}

func (c *Cache) SetEvictionStrategy(algo EvictionAlgo) {
	c.evictionAlgo = algo
}

func initCache(e EvictionAlgo) *Cache {
	return &Cache{
		storage:      sync.Map{},
		evictionAlgo: e,
		capacity:     0,
		maxCapacity:  500,
	}
}
