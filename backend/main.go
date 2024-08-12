package main

import (
	"log"
	"net/http"
	"time"

	"github.com/rs/cors"
)

type CacheItem struct {
	Value      string
	Expiration time.Time
}

type LRUCache struct {
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
	if _, ok := c.items[key]; ok {
		c.moveToFront(key)
	} else {
		if len(c.items) >= c.capacity {
			oldest := c.keys[len(c.keys)-1]
			delete(c.items, oldest)
			c.removeKey(oldest)
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

func main() {
	cache = NewLRUCache(1024)

	mux := http.NewServeMux()
	mux.HandleFunc("/get", getHandler)
	mux.HandleFunc("/set", setHandler)

	// Create a CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Allow requests from the React app
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
	})

	// Wrap the server with CORS middleware
	handler := c.Handler(mux)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
