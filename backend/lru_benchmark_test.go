package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func BenchmarkLRUCache(b *testing.B) {
	cacheSizes := []int{100, 1000, 10000}
	operationCounts := []int{1000, 10000, 100000}

	for _, size := range cacheSizes {
		for _, count := range operationCounts {
			b.Run(fmt.Sprintf("Size_%d_Ops_%d", size, count), func(b *testing.B) {
				cache := NewLRUCache(size)
				operation := make([]bool, count)
				keys := make([]string, count)
				values := make([]string, count)

				for i := 0; i < count; i++ {
					operation[i] = rand.Int()%2 == 0
					keys[i] = randomString(10)
					values[i] = randomString(20)
				}

				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					// Benchmark Set operations
					for j := 0; j < count; j++ {
						if operation[j] {
							cache.Set(keys[j], values[j], 10*time.Millisecond)
						} else {
							cache.Get(keys[rand.Intn(count)])
						}
					}
				}
			})
		}
	}
}

func BenchmarkLRUCacheParallel(b *testing.B) {
	cache := NewLRUCache(10000)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := randomString(10)
			value := randomString(20)
			cache.Set(key, value, time.Minute)
			cache.Get(key)
		}
	})
}
