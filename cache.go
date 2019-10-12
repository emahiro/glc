/*
Package cache is provides the local cache which is stored in memoroy or file.

Example:
	mc := NewMemoryCache(cache.DefaultMemoryCacheExpires)

	// Set
	if err := mc.Set("cacheKey", []byte('hoge')); err != nil {
		log.Fatal(err)
	}

	// Get
	data := mc.Get("cacheKey")
*/
package cache

import (
	"fmt"
	"sync"
	"time"
)

// DefaultMemoryCacheExpires is 60 seconds
const DefaultMemoryCacheExpires = 60 * time.Second

// MemoryCache is cache data in memory which has expiration date.
type MemoryCache struct {
	data    map[string][]byte
	expires int64
	m       sync.RWMutex
}

// Get returns a item or nil.
// If cache in local is nil or expiration date of the cache you want to retrive is earlier, you can't retrive cache.
func (c *MemoryCache) Get(key string) []byte {
	c.m.RLock()
	defer c.m.RUnlock()

	now := time.Now().Unix()
	if c == nil || c.expires < now {
		return nil
	}

	if data, ok := c.data[key]; ok && len(data) != 0 {
		return data
	}
	return nil
}

// Set add a new data for cache with a new key or replace an exist key.
func (c *MemoryCache) Set(key string, src []byte) error {
	if c.data == nil {
		return fmt.Errorf("error: nil map access")
	}

	if len(src) == 0 {
		return fmt.Errorf("error: set no data")
	}

	c.m.Lock()
	defer c.m.Unlock()

	c.data[key] = src
	return nil
}

// NewMemoryCache creates a new MemoryCache for given a its expires.
// If exp is 0, you will use the default cache expiration.
// The default cache expiration is 60 seconds.
func NewMemoryCache(exp time.Duration) *MemoryCache {
	if exp == 0 {
		exp = DefaultMemoryCacheExpires
	}
	return &MemoryCache{data: map[string][]byte{}, expires: time.Now().Add(exp).Unix()}
}
