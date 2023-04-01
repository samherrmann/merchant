// Package cache handles the local caching of store data.
package cache

import (
	"errors"
	"os"

	"github.com/samherrmann/merchant/osutil"
)

const (
	AppName    = "merchant"
	dbFilename = "bolt.db"
)

var (
	ErrExist    = errors.New("already exists")
	ErrNotExist = errors.New("does not exist")
)

type Cache interface {
	Products() ProductCache
}

// New returns a new cache.
func New() (Cache, error) {
	dbOpener, err := newDBOpener()
	if err != nil {
		return nil, err
	}
	cache := &cache{
		products: NewProductCache(dbOpener),
	}
	return cache, nil
}

type cache struct {
	products ProductCache
}

func (c *cache) Products() ProductCache {
	return c.products
}

// Clear removes the cache directory.
func Clear() error {
	dir, err := directory()
	if err != nil {
		return err
	}
	return os.RemoveAll(dir)
}

// directory returns the path to the cache directory. If the directory does not
// exist, then directory will create it.
func directory() (string, error) {
	cacheRootDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return osutil.MakeUserDir(cacheRootDir, AppName)
}
