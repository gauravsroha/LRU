package main

import (
	"sync"
	"time"
)

type CacheItem struct {
	Value      string
	Expiration time.Time
}

type LRUCache struct {
	sync.Mutex
	capacity int
	items    map[string]CacheItem
	keys     []string
}

var cache *LRUCache

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		items:    make(map[string]CacheItem),
		keys:     make([]string, 0), //initialise the key of capacity
	}
}

func (c *LRUCache) Get(key string) (string, bool) {
	c.Lock()
	defer c.Unlock()
	if item, ok := c.items[key]; ok {
		if time.Now().After(item.Expiration) {
			delete(c.items, key)
			c.removeKey(key)
			return "", false
		}
		c.moveToFront(key)
		return item.Value, true
	}
	return "", false
}

func (c *LRUCache) Set(key, value string, expiration time.Duration) {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.items[key]; ok {
		c.moveToFront(key)
	} else {
		if len(c.items) >= c.capacity {
			oldest := c.keys[len(c.keys)-1]
			delete(c.items, oldest)
			c.keys = c.keys[:len(c.keys)-1]
		}
		c.keys = append([]string{key}, c.keys...)
	}
	c.items[key] = CacheItem{
		Value:      value,
		Expiration: time.Now().Add(expiration),
	}
}

func (c *LRUCache) moveToFront(key string) {
	for i, k := range c.keys {
		if k == key {
			c.keys = append(c.keys[:i], c.keys[i+1:]...)
			break
		}
	}
	c.keys = append([]string{key}, c.keys...)
}

func (c *LRUCache) removeKey(key string) {
	for i, k := range c.keys {
		if k == key {
			c.keys = append(c.keys[:i], c.keys[i+1:]...)
			break
		}
	}
}
