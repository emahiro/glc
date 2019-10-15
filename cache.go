/*
Package cache is provides the local cache which is stored in memoroy or file.
This package creates `tmp` directory for file cache, when you provide UseFileCache true.

Example:
	mc := NewMemoryCache(time.Now().Add(cache.DefaultMemoryCacheExpires*time.Second))

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

var (
	// UseFileCache is flag whitch control file cache usage
	UseFileCache = false
)

// MemoryCache is cache data in memory which has duration.
type MemoryCache struct {
	item map[string]*Item
	d    time.Duration
	m    sync.RWMutex
}

// Item has cache item and expiration.
type Item struct {
	data []byte
	exp  int64
}

func init() {
	if !UseFileCache {
		return
	}

	if _, err := os.Stat(fileCacheDir); os.IsNotExist(err) {
		if err := os.Mkdir(fileCacheDir, os.ModePerm); err != nil {
			panic(err)
		}
	}
}

// Get returns a item or nil.
// If cache in local is nil or expiration date of the cache you want to retrive is earlier, you can't retrive cache.
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

// NewMemoryCache creates a new MemoryCache for given a its expires as time.Time.
// If exp is 0, you will use the default cache expiration.
// The default cache expiration is 60 seconds.
func NewMemoryCache(d time.Duration) *MemoryCache {
	return &MemoryCache{item: make(map[string]*Item), d: d}
}

// FileCache is cache data in local file.
type FileCache struct {
	m sync.RWMutex
}

// Get returns a data or nil.
// If cache in local file is nil or you set key that does not hit the cache, you can not retrive cache.
func (c *FileCache) Get(key string) []byte {
	c.m.RLock()
	defer c.m.RUnlock()

	fp := filepath.Join(".", fileCacheDir, fmt.Sprintf("%s.cache", key))
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil
	}

	return b
}

// Set create a new file which is witten data as []byte.
// When you set a new cache data, create a `{{ $key }}.cache` file,
// and if cache file is exist, overwrite it.
func (c *FileCache) Set(key string, src []byte) error {
	if len(src) == 0 {
		return fmt.Errorf("error: set no data")
	}

	c.m.Lock()
	defer c.m.Unlock()

	fp := filepath.Join(".", fileCacheDir, fmt.Sprintf("%s.cache", key))
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		if _, err := os.Create(fp); err != nil {
			return fmt.Errorf("set cache error. err:%v", err)
		}
	}

	if err := ioutil.WriteFile(fp, src, os.ModePerm); err != nil {
		return fmt.Errorf("set cache error. err: %v", err)
	}

	return nil
}
