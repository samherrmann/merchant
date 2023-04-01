package cache

import (
	goshopify "github.com/bold-commerce/go-shopify/v3"
	bolt "go.etcd.io/bbolt"
)

type ProductCache interface {
	Update(p ...goshopify.Product) error
	GetByID(id int64) (*goshopify.Product, error)
	GetByTitle(title string) (*goshopify.Product, error)
	GetByHandle(handle string) (*goshopify.Product, error)
	List() ([]goshopify.Product, error)
}

func NewProductCache(o DBOpener) ProductCache {
	return &productCache{dbOpener: o}
}

type productCache struct {
	dbOpener DBOpener
}

func (cache *productCache) Update(products ...goshopify.Product) error {
	return cache.update(func(b *ProductBuckets) error {
		return b.Update(products...)
	})
}

func (cache *productCache) GetByID(id int64) (p *goshopify.Product, err error) {
	err = cache.view(func(b *ProductBuckets) error {
		p, err = b.GetByID(id)
		return err
	})
	return p, err
}

func (cache *productCache) GetByTitle(title string) (p *goshopify.Product, err error) {
	err = cache.view(func(buckets *ProductBuckets) error {
		p, err = buckets.GetByTitle(title)
		return err
	})
	return p, err
}

func (cache *productCache) GetByHandle(handle string) (p *goshopify.Product, err error) {
	err = cache.view(func(buckets *ProductBuckets) error {
		p, err = buckets.GetByHandle(handle)
		return err
	})
	return p, err
}

func (cache *productCache) List() (p []goshopify.Product, err error) {
	err = cache.view(func(buckets *ProductBuckets) error {
		p, err = buckets.List()
		return err
	})
	return p, err
}

func (cache *productCache) view(fn func(b *ProductBuckets) error) error {
	db, err := cache.dbOpener.Open()
	if err != nil {
		return err
	}
	defer db.Close()
	return db.View(func(tx *bolt.Tx) error {
		buckets, err := NewProductBuckets(tx)
		if err != nil {
			return err
		}
		return fn(buckets)
	})
}

func (cache *productCache) update(fn func(b *ProductBuckets) error) error {
	db, err := cache.dbOpener.Open()
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Update(func(tx *bolt.Tx) error {
		buckets, err := NewProductBuckets(tx)
		if err != nil {
			return err
		}
		return fn(buckets)
	})
}
