package cache

import (
	"fmt"
	"sync"
	"time"
)

type Cache interface {
	Get(key string) []byte
	Set(key string, src []byte) error
}

type LocalCache struct {
	Data    map[string][]byte
	Expires int64
	m       sync.RWMutex
}

// Get
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

// Set
func (c *LocalCache) Set(key string, src []byte) error {
	c.m.Lock()
	defer c.m.Unlock()

	if len(src) == 0 {
		return fmt.Errorf("no set data")
	}

	c.Data[key] = src
	return nil
}

// NewLocalCache
func NewLocalCache(exp int64) Cache {
	return &LocalCache{Data: map[string][]byte{}, Expires: exp}
}
