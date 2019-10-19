/*
Package glc is provides the local cache which is stored in memoroy or file.
This package creates `tmp` directory in case of using file cache.

Example:
	mc := glc.NewMemoryCache(glc.DefaultMemoryCacheExpires)

	// Set
	if err := mc.Set("cacheKey", []byte('hoge')); err != nil {
		log.Fatal(err)
	}

	// Get
	data := mc.Get("cacheKey")
*/
package glc

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	// DefaultMemoryCacheExpires is 60 seconds
	DefaultMemoryCacheExpires = 60 * time.Second
	fileCacheDir              = "tmp"
)

// MemoryCache is cache data in memory which has duration.
type MemoryCache struct {
	item map[string]*Item
	d    time.Duration
	m    sync.RWMutex
}

// Item has cache item and expiration field.
type Item struct {
	data []byte
	exp  int64
}

// Get returns a item or nil.
// If cache in local is nil or expiration date of the cache is earlier, this returns nil.
func (c *MemoryCache) Get(key string) []byte {
	c.m.RLock()
	defer c.m.RUnlock()

	item, ok := c.item[key]
	if !ok || item == nil {
		return nil
	}
	if len(item.data) == 0 || item.exp < time.Now().UnixNano() {
		return nil
	}

	return item.data
}

// Set add a new data for cache with a new key or replace an exist key.
func (c *MemoryCache) Set(key string, src []byte) error {
	if c.item == nil {
		return fmt.Errorf("error: nil map access")
	}

	if len(src) == 0 {
		return fmt.Errorf("error: set no data")
	}

	c.m.Lock()
	defer c.m.Unlock()

	c.item[key] = &Item{
		data: src,
		exp:  time.Now().Add(c.d).UnixNano(),
	}
	return nil
}

// NewMemoryCache creates a new MemoryCache for given a its expires as time.Duration.
func NewMemoryCache(d time.Duration) *MemoryCache {
	return &MemoryCache{item: make(map[string]*Item), d: d}
}

// FileCache is cache data in local file.
type FileCache struct {
	path string
	m    sync.RWMutex
}

// Get returns a data or nil.
// If cache in local file is nil or is not setted key, this returns nil.
func (c *FileCache) Get(key string) []byte {
	c.m.RLock()
	defer c.m.RUnlock()

	fp := filepath.Join(c.path, key) + ".cache"
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil
	}

	return b
}

// Set create a new file which is witten data as []byte.
// When it sets a new cache data, create a `{{ $key }}.cache` file,
// and overwrite it if cache file is exist.
func (c *FileCache) Set(key string, src []byte) error {
	if len(src) == 0 {
		return fmt.Errorf("error: set no data")
	}

	c.m.Lock()
	defer c.m.Unlock()

	fp := filepath.Join(c.path, key) + ".cache"
	if err := ioutil.WriteFile(fp, src, os.ModePerm); err != nil {
		return fmt.Errorf("set cache error. err: %v", err)
	}

	return nil
}

// NewFileCache returns FileCache pointer.
// If there is no temp directory named prefix, create temp directory for storing cache.
func NewFileCache(prefix string) (*FileCache, error) {
	path, err := ioutil.TempDir("", prefix)
	if err != nil {
		return nil, err
	}

	return &FileCache{path: path}, nil
}
