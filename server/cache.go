package server

import (
	"container/list"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type value struct {
	key        string
	data       string
	accessed   int32
	ttl        time.Duration
	elementPtr *list.Element
}

type EvictionAlgo interface {
	evict(*Cache) error
	ttl() time.Duration
}

type Eviction struct {
	EvictionAlgo
}

type Cache struct {
	store        sync.Map
	list         *list.List
	evictionAlgo EvictionAlgo
	capacity     int
	maxCapacity  int
}

// Returns a new cache instance with LRU by default
func NewCache() *Cache {
	var l list.List
	l.Init()
	c := &Cache{
		store:        sync.Map{},
		capacity:     0,
		maxCapacity:  25,
		evictionAlgo: &LRU{},
		list:         &l,
	}
	return c
}

func (c *Cache) Get(key string) (string, error) {
	val, ok := c.store.Load(key)
	if ok {
		v, ok := val.(value)
		if ok {
			c.list.MoveToBack(v.elementPtr)
			v.accessed++
			c.store.Store(key, v)
			return v.data, nil
		}
		return "", errors.New("cache: type assertion failed")
	}
	return "", errors.New("cache: key not found")
}

// Puts a key-value pair into cache, returns false if key already exists
func (c *Cache) Set(key string, val string) error {
	if c.capacity >= c.maxCapacity {
		if err := c.evictionAlgo.evict(c); err != nil {
			return err
		}
	}

	if _, ok := c.store.Load(key); ok {
		return fmt.Errorf("cache: key already exists")
	}

	var v value

	v.data = val
	v.accessed = 1
	v.ttl = c.evictionAlgo.ttl()
	v.key = key
	v.elementPtr = c.list.PushBack(v)

	c.store.Store(
		key,
		v,
	)
	c.capacity++

	if v.ttl != 0 {
		go func() {
			time.AfterFunc(v.ttl, func() {
				log.Println("ttl evicted", key)
				c.store.Delete(key)
				c.list.Remove(v.elementPtr)
				c.capacity--
			})
		}()
	}

	return nil
}

func (c *Cache) GetCacheData() map[string]value {
	m := make(map[string]value)
	c.store.Range(func(k any, v any) bool {
		m[k.(string)] = v.(value)
		return true
	})
	return m
}

func (c *Cache) SetEvictionStrategy(algo EvictionAlgo) {
	c.evictionAlgo = algo
}
