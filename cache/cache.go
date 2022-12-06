// A simple Put/Get wrapper around bbolt [1].
//
// [1]: https://github.com/etcd-io/bbolt
package cache

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	bolt "go.etcd.io/bbolt"
)

type (
	Cache interface {
		Put(key string, value []byte) error
		Get(key string) ([]byte, error)
	}

	BoltCache struct {
		cacheDirectory string
		cacheName      string
		lifetime       time.Duration
		db             *bolt.DB
	}

	cacheValue struct {
		value []byte
		ts    time.Time
	}

	cacheValueJSON struct {
		Value []byte
		Ts    []byte
	}
)

func (cv cacheValue) MarshalJSON() ([]byte, error) {
	store := cacheValueJSON{
		Value: cv.value,
		Ts:    []byte(cv.ts.Format(time.RFC3339)),
	}

	v, err := json.Marshal(store)
	return v, err
}

func (cv *cacheValue) UnmarshalJSON(data []byte) error {
	var store cacheValueJSON
	err := json.Unmarshal(data, &store)
	if err != nil {
		return err
	}

	ts, err := time.Parse(time.RFC3339, string(store.Ts))
	if err != nil {
		return err
	}

	cv.value = store.Value
	cv.ts = ts

	return nil
}

func NewCache(appName, cacheName string) *BoltCache {
	cacheDirectory := filepath.Join(xdg.CacheHome, appName)
	return &BoltCache{
		cacheDirectory: cacheDirectory,
		cacheName:      cacheName,
	}
}

func (cache *BoltCache) WithLifetime(lifetime time.Duration) *BoltCache {
	cache.lifetime = lifetime
	return cache
}

func (cache *BoltCache) WithCacheDirectory(dir string) *BoltCache {
	cache.cacheDirectory = dir
	return cache
}

func (cache *BoltCache) Start() error {
	err := ensureDir(cache.cacheDirectory, 0755)
	if err != nil {
		return err
	}

	db, err := bolt.Open(filepath.Join(cache.cacheDirectory, "cache.db"), 0600, nil)
	if err != nil {
		return err
	}
	cache.db = db

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(cache.cacheName))
		return err
	})
	return err
}

func (cache *BoltCache) Get(key string) ([]byte, error) {
	var data []byte
	var cv cacheValue

	// We don't check the return value here because this
	// always returns successfully.
	//nolint:errcheck
	cache.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(cache.cacheName))
		data = b.Get([]byte(key))
		return nil
	})

	if data == nil {
		return nil, nil
	}

	if err := json.Unmarshal(data, &cv); err != nil {
		return nil, err
	}

	// Return nil if cache value has expired.
	if cache.lifetime > 0 && time.Since(cv.ts) > cache.lifetime {
		return nil, nil
	}

	return cv.value, nil
}

func (cache *BoltCache) Put(key string, value []byte) error {
	data, err := json.Marshal(cacheValue{
		value: value,
		ts:    time.Now(),
	})
	if err != nil {
		return err
	}

	err = cache.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(cache.cacheName))
		err := b.Put([]byte(key), data)
		return err
	})

	return err
}

// via https://stackoverflow.com/a/56600630/147356
func ensureDir(dirName string, mode os.FileMode) error {
	err := os.Mkdir(dirName, mode)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		// check that the existing path is a directory
		info, err := os.Stat(dirName)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return errors.New("path exists but is not a directory")
		}
		return nil
	}
	return err
}
