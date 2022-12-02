package cache

import (
	"errors"
	"os"
	"path/filepath"

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
		db             *bolt.DB
	}
)

func NewCache(appName, cacheName string) *BoltCache {
	cacheDirectory := filepath.Join(xdg.CacheHome, appName)
	return &BoltCache{
		cacheDirectory: cacheDirectory,
		cacheName:      cacheName,
	}
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
	return nil
}

func (cache *BoltCache) Get(key string) ([]byte, error) {
	var value []byte

	cache.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(cache.cacheName))
		value = b.Get([]byte(key))
		return nil
	})

	return value, nil
}

func (cache *BoltCache) Put(key string, value []byte) error {
	err := cache.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(cache.cacheName))
		err := b.Put([]byte(key), []byte(value))
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
