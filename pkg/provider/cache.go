package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/gofrs/flock"
	"github.com/jellydator/ttlcache/v2"
)

type cacheConfig struct {
	CacheFile string
	LockFile  string
}

type cacheStore struct {
	cache     *ttlcache.Cache
	cacheFile string
	lockFile  string
}

type cacheData struct {
	Data []CarbonIntensity `json:"data"`
	TTL  time.Time         `json:"ttl"`
}

func NewCacheStore(config cacheConfig) (*cacheStore, error) {
	var lockFile string

	if config.CacheFile != "" {
		lockFile = config.CacheFile + ".lock"
	}

	store := &cacheStore{
		cache:     ttlcache.NewCache(),
		cacheFile: config.CacheFile,
		lockFile:  lockFile,
	}

	return store, nil
}

func (c *cacheStore) getCacheData(ctx context.Context, region string) ([]CarbonIntensity, error) {
	var result *cacheData

	if c.cacheFile == "" {
		raw, err := c.cache.Get(region)
		if errors.Is(ttlcache.ErrNotFound, err) {
			// Cache miss so return nil.
			return nil, nil
		} else if err != nil {
			return nil, err
		}

		item, ok := raw.(*cacheData)
		if !ok {
			return nil, fmt.Errorf("cannot convert %#v to %T", raw, []CarbonIntensity{})
		}
		result = item
	} else {
		cache, err := c.loadCache(ctx)
		if err != nil {
			return nil, err
		}

		item, ok := cache[region]
		if !ok {
			// Cache miss so return nil.
			return nil, nil
		}
		result = item
	}

	if result.TTL.Before(time.Now()) {
		// Item has expired.
		return nil, nil
	}

	return result.Data, nil
}

func (c *cacheStore) setCacheData(ctx context.Context, region string, data []CarbonIntensity, ttl time.Time) error {
	item := &cacheData{
		Data: data,
		TTL:  ttl,
	}

	if c.cacheFile == "" {
		c.cache.Set(region, item)
	} else {
		err := c.saveCache(ctx, region, item)
		if err != nil {
			return nil
		}
	}

	return nil
}

func (c *cacheStore) loadCache(ctx context.Context) (map[string]*cacheData, error) {
	cache := make(map[string]*cacheData, 0)

	// Ensure the directory exists as it is required for file locking.
	err := os.MkdirAll(filepath.Dir(c.cacheFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}

	lockCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Get a shared file lock for reading.
	fileLock := flock.New(c.lockFile)
	locked, err := fileLock.TryRLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer fileLock.Unlock()
	}
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(c.cacheFile)
	if _, ok := err.(*fs.PathError); ok {
		return cache, nil
	} else if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &cache)
	if err != nil {
		return nil, err
	}

	return cache, nil
}

func (c *cacheStore) saveCache(ctx context.Context, region string, item *cacheData) error {
	cache, err := c.loadCache(ctx)
	if err != nil {
		return err
	}
	cache[region] = item

	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	lockCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Get an exclusive file lock for writing.
	fileLock := flock.New(c.lockFile)
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer fileLock.Unlock()
	}
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.cacheFile, data, 0644)
	if err != nil {
		return nil
	}

	return nil
}
