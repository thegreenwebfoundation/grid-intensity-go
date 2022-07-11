package watttime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofrs/flock"
	"github.com/jellydator/ttlcache/v2"
)

func (a *ApiClient) getCacheData(ctx context.Context, region string) (*IndexData, error) {
	var result *CacheData

	if a.cacheFile == "" {
		raw, err := a.cache.Get(region)
		if errors.Is(ttlcache.ErrNotFound, err) {
			// Cache miss so return nil.
			return nil, nil
		} else if err != nil {
			return nil, err
		}

		item, ok := raw.(*CacheData)
		if !ok {
			return nil, fmt.Errorf("cannot convert %#v to %T", raw, &CacheData{})
		}
		result = item
	} else {
		cache, err := a.loadCache(ctx)
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

func (a *ApiClient) loadCache(ctx context.Context) (map[string]*CacheData, error) {
	cache := make(map[string]*CacheData, 0)

	// Ensure the directory exists as it is required for file locking.
	err := os.MkdirAll(filepath.Dir(a.cacheFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}

	lockCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Get a shared file lock for reading.
	fileLock := flock.New(a.lockFile)
	locked, err := fileLock.TryRLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer fileLock.Unlock()
	}
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(a.cacheFile)
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

func (a *ApiClient) saveCache(ctx context.Context, region string, item *CacheData) error {
	cache, err := a.loadCache(ctx)
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
	fileLock := flock.New(a.lockFile)
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer fileLock.Unlock()
	}
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(a.cacheFile, data, 0644)
	if err != nil {
		return nil
	}

	return nil
}

func (a *ApiClient) setCacheData(ctx context.Context, region string, result *IndexData) error {
	freq, err := strconv.ParseInt(result.Freq, 0, 64)
	if err != nil {
		return err
	}

	ttl := result.PointTime.Add(time.Duration(freq) * time.Second)
	if ttl.Before(time.Now()) {
		// The TTL calculated from the point time is in the past. So reset the
		// TTL using the current time plus the frequency provided by the API.
		// UTC is used to match the WattTime API.
		ttl = time.Now().UTC().Add(time.Duration(freq) * time.Second)
	}
	item := &CacheData{
		Data: result,
		TTL:  ttl,
	}

	if a.cacheFile == "" {
		a.cache.Set(region, item)
	} else {
		err := a.saveCache(ctx, region, item)
		if err != nil {
			return nil
		}
	}

	return nil
}
