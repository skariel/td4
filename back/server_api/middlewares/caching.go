package middlewares

import (
	"log"
	"time"

	cache "github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"
)

// NewMemoryCache a memory cache
func NewMemoryCache(capacity int, ttl time.Duration) *cache.Client {
	memcached, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(capacity),
	)
	if err != nil {
		log.Fatal(err)
	}

	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(memcached),
		cache.ClientWithTTL(ttl),
		cache.ClientWithRefreshKey("uncached"),
	)
	if err != nil {
		log.Fatal(err)
	}

	return cacheClient
}
