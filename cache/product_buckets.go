package cache

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/samherrmann/merchant/cache/bkeys"
	bolt "go.etcd.io/bbolt"
)

func NewProductBuckets(tx *bolt.Tx) (*ProductBuckets, error) {
	return &ProductBuckets{
		tx:       tx,
		products: tx.Bucket([]byte(bkeys.Products)),
		handles:  tx.Bucket([]byte(bkeys.ProductHandles)),
		titles:   tx.Bucket([]byte(bkeys.ProductTitles)),
	}, nil
}

// ProductBuckets is a collection of Bolt Buckets to store products.
type ProductBuckets struct {
	tx       *bolt.Tx
	products *bolt.Bucket
	titles   *bolt.Bucket
	handles  *bolt.Bucket
}

func (b *ProductBuckets) GetByID(id int64) (*goshopify.Product, error) {
	if b.products == nil {
		return nil, ErrNotExist
	}
	k := int64ToBytes(id)
	v := b.products.Get(k)
	if v == nil {
		return nil, ErrNotExist
	}
	p := &goshopify.Product{}
	err := json.Unmarshal(v, p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (b *ProductBuckets) GetByHandle(handle string) (*goshopify.Product, error) {
	return b.getBySecondaryKey(b.handles, []byte(handle))
}

func (b *ProductBuckets) GetByTitle(title string) (*goshopify.Product, error) {
	return b.getBySecondaryKey(b.titles, []byte(title))
}

func (b *ProductBuckets) Update(products ...goshopify.Product) error {
	for _, p := range products {
		// Insert in handle bucket.
		if p.Handle != "" {
			var err error
			b.handles, err = b.tx.CreateBucketIfNotExists([]byte(bkeys.ProductHandles))
			if err != nil {
				return err
			}
			k := []byte(p.Handle)
			v := int64ToBytes(p.ID)
			if err := setOnce(b.handles, k, v); err != nil {
				return fmt.Errorf("product handle %q: %w", p.Handle, err)
			}
		}
		// Insert in title bucket.
		if p.Title != "" {
			var err error
			b.titles, err = b.tx.CreateBucketIfNotExists([]byte(bkeys.ProductTitles))
			if err != nil {
				return err
			}
			k := []byte(p.Title)
			v := int64ToBytes(p.ID)
			if err := setOnce(b.titles, k, v); err != nil {
				return fmt.Errorf("product title %q: %w", p.Title, err)
			}
		}
		// Insert in main bucket:
		{
			var err error
			b.products, err = b.tx.CreateBucketIfNotExists([]byte(bkeys.Products))
			if err != nil {
				return err
			}
			k := int64ToBytes(p.ID)
			v, err := json.Marshal(p)
			if err != nil {
				return err
			}
			if err := b.products.Put(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *ProductBuckets) List() ([]goshopify.Product, error) {
	if b.products == nil {
		return nil, nil
	}
	var products []goshopify.Product
	err := b.products.ForEach(func(k, v []byte) error {
		p := goshopify.Product{}
		if err := json.Unmarshal(v, &p); err != nil {
			return err
		}
		products = append(products, p)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return products, nil
}

// getProductBySecondaryKey returns the product from the primary bucket given a
// secondary key. The value associated with the secondary key is expected to be
// the key of the product in the primary bucket.
func (b *ProductBuckets) getBySecondaryKey(bucket *bolt.Bucket, key []byte) (*goshopify.Product, error) {
	if bucket == nil {
		return nil, ErrNotExist
	}
	if b.products == nil {
		return nil, ErrNotExist
	}

	primaryKey := bucket.Get(key)
	if primaryKey == nil {
		return nil, ErrNotExist
	}
	v := b.products.Get(primaryKey)
	if v == nil {
		return nil, ErrNotExist
	}
	p := &goshopify.Product{}
	if err := json.Unmarshal(v, p); err != nil {
		return nil, err
	}
	return p, nil
}

// setOnce sets the value for a key in the bucket. If the key already exists,
// then the new value must be the same as the one that is already in the bucket.
// If the values do not match then setOnce returns an error.
func setOnce(bucket *bolt.Bucket, k []byte, v []byte) error {
	dbv := bucket.Get(k)
	if dbv != nil {
		if bytes.Equal(dbv, v) {
			return nil
		}
		return fmt.Errorf("cannot set key %q to %q: already set to %q", k, v, dbv)
	}
	return bucket.Put(k, v)
}

// int64ToBytes returns the byte encoding of v.
func int64ToBytes(v int64) []byte {
	return []byte(strconv.FormatInt(v, 10))
}
