package cache

import (
	"fmt"
	"sync"
	"time"
)

var (
	// DefaultLocalCacheExpires is 60 seconds
	DefaultLocalCacheExpires int64 = 60
)

// Cache is the interface that contains Get and Set method.
// Get returns cached item as []byte.
// Set add or replace a new item as []byte with a new key or an exist key.
type Cache interface {
	Get(key string) []byte
	Set(key string, src []byte) error
}

// LocalCache is cache data in local which has expiration date.
type LocalCache struct {
	Data    map[string][]byte
	Expires int64
	m       sync.RWMutex
}

// Get returns a item or nil.
// If cache in local is nil or expiration date of the cache you want to retrive is earlier, you can't retrive cache.
func (c *LocalCache) Get(key string) []byte {
	c.m.RLock()
	defer c.m.RUnlock()

	now := time.Now().Unix()
	if c == nil || c.Expires < now {
		return nil
	}

	if data, ok := c.Data[key]; ok && len(data) != 0 {
		return data
	}
	return nil
}

// Set add a new data for cache with a new key or replace an exist key.
func (c *LocalCache) Set(key string, src []byte) error {
	if c.Data == nil {
		return fmt.Errorf("error: nil map access")
	}

	if len(src) == 0 {
		return fmt.Errorf("error: set no data")
	}

	c.m.Lock()
	defer c.m.Unlock()

	c.Data[key] = src
	return nil
}

// NewLocalCache creates a new LocalCache for given a its expires.
// If exp is 0, you will use the default cache expiration.
// The default cache expiration is 60 seconds.
func NewLocalCache(exp int64) Cache {
	if exp == 0 {
		exp = DefaultLocalCacheExpires
	}
	return &LocalCache{Data: map[string][]byte{}, Expires: exp}
}
