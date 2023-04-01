package cache

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/samherrmann/merchant/cache/bkeys"
	bolt "go.etcd.io/bbolt"
)

func TestProductBuckets_GetByID(t *testing.T) {
	t.Run("returns ErrNotExist if bucket doesn't exist", func(t *testing.T) {
		db := newTestDB(t)
		db.View(func(tx *bolt.Tx) error {
			buckets, err := NewProductBuckets(tx)
			if err != nil {
				t.Fatal(err)
			}
			want := ErrNotExist
			_, got := buckets.GetByID(123)
			if !errors.Is(got, want) {
				t.Fatalf("got %v, want %v", got, want)
			}
			return nil
		})
	})

	t.Run("returns ErrNotExist if product doesn't exist", func(t *testing.T) {
		db := newTestDB(t)

		// Add a product so that the bucket gets created.
		db.Update(func(tx *bolt.Tx) error {
			buckets, err := NewProductBuckets(tx)
			if err != nil {
				t.Fatal(err)
			}
			if err := buckets.Update(goshopify.Product{}); err != nil {
				t.Fatal(err)
			}
			return nil
		})

		db.View(func(tx *bolt.Tx) error {
			buckets, err := NewProductBuckets(tx)
			if err != nil {
				t.Fatal(err)
			}
			want := ErrNotExist
			_, got := buckets.GetByID(1)
			if !errors.Is(got, want) {
				t.Fatalf("got %v, want %v", got, want)
			}
			return nil
		})
	})

	t.Run("returns existing product", func(t *testing.T) {
		p := goshopify.Product{ID: 123}

		db := newTestDB(t)
		db.Update(func(tx *bolt.Tx) error {
			buckets, err := NewProductBuckets(tx)
			if err != nil {
				t.Fatal(err)
			}
			if err := buckets.Update(p); err != nil {
				t.Fatal(err)
			}
			return nil
		})

		db.View(func(tx *bolt.Tx) error {
			buckets, err := NewProductBuckets(tx)
			if err != nil {
				t.Fatal(err)
			}
			got, err := buckets.GetByID(123)
			if err != nil {
				t.Fatal(err)
			}
			want := &p
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("\ngot:\n%v\nwant:\n%v\n", got, want)
			}
			return nil
		})
	})
}

func TestProductBuckets_getBySecondaryKey(t *testing.T) {
	t.Run("returns ErrNotExist if main bucket doesn't exist", func(t *testing.T) {
		db := newTestDB(t)
		db.Update(func(tx *bolt.Tx) error {
			buckets, err := NewProductBuckets(tx)
			if err != nil {
				t.Fatal(err)
			}
			want := ErrNotExist
			_, got := buckets.getBySecondaryKey(buckets.handles, []byte("foo"))
			if !errors.Is(got, want) {
				t.Fatalf("got %v, want %v", got, want)
			}
			return nil
		})
	})

	t.Run("returns ErrNotExist if secondary bucket doesn't exist", func(t *testing.T) {
		db := newTestDB(t)

		// Add a product so that the main bucket gets created.
		db.Update(func(tx *bolt.Tx) error {
			buckets, err := NewProductBuckets(tx)
			if err != nil {
				t.Fatal(err)
			}
			p := goshopify.Product{}
			if err := buckets.Update(p); err != nil {
				t.Fatal(err)
			}
			return nil
		})

		db.View(func(tx *bolt.Tx) error {
			buckets, err := NewProductBuckets(tx)
			if err != nil {
				t.Fatal(err)
			}
			want := ErrNotExist
			_, got := buckets.getBySecondaryKey(buckets.handles, []byte("foo"))
			if !errors.Is(got, want) {
				t.Fatalf("got %v, want %v", got, want)
			}
			return nil
		})
	})

	t.Run("returns product", func(t *testing.T) {
		p := goshopify.Product{ID: 123, Handle: "foo"}

		db := newTestDB(t)
		db.Update(func(tx *bolt.Tx) error {
			buckets, err := NewProductBuckets(tx)
			if err != nil {
				t.Fatal(err)
			}
			if err := buckets.Update(p); err != nil {
				t.Fatal(err)
			}
			return nil
		})

		db.View(func(tx *bolt.Tx) error {
			buckets, err := NewProductBuckets(tx)
			if err != nil {
				t.Fatal(err)
			}
			want := &p
			got, err := buckets.getBySecondaryKey(buckets.handles, []byte("foo"))
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("\ngot:\n%v\nwant:\n%v\n", got, want)
			}
			return nil
		})
	})

	t.Run("returns ErrNotExist if not found in secondary bucket", func(t *testing.T) {
		p := goshopify.Product{ID: 123}

		db := newTestDB(t)
		db.Update(func(tx *bolt.Tx) error {
			buckets, err := NewProductBuckets(tx)
			if err != nil {
				t.Fatal(err)
			}
			if err := buckets.Update(p); err != nil {
				t.Fatal(err)
			}
			return nil
		})

		db.View(func(tx *bolt.Tx) error {
			buckets, err := NewProductBuckets(tx)
			if err != nil {
				t.Fatal(err)
			}
			want := ErrNotExist
			_, got := buckets.getBySecondaryKey(buckets.handles, []byte("foo"))
			if !errors.Is(got, want) {
				t.Fatalf("got %v, want %v", got, want)
			}
			return nil
		})
	})
}

func TestProductBuckets_List(t *testing.T) {

	t.Run("returns empty slice if primary bucket does not exist", func(t *testing.T) {
		db := newTestDB(t)
		db.View(func(tx *bolt.Tx) error {
			buckets, err := NewProductBuckets(tx)
			if err != nil {
				t.Fatal(err)
			}
			products, err := buckets.List()
			if err != nil {
				t.Fatal(err)
			}
			want := 0
			got := len(products)
			if got != want {
				t.Fatalf("got %v, want %v", got, want)
			}
			return nil
		})
	})

	t.Run("returns products", func(t *testing.T) {
		products := []goshopify.Product{
			{ID: 1},
			{ID: 2},
			{ID: 3},
		}

		db := newTestDB(t)
		db.Update(func(tx *bolt.Tx) error {
			buckets, err := NewProductBuckets(tx)
			if err != nil {
				t.Fatal(err)
			}
			if err := buckets.Update(products...); err != nil {
				t.Fatal(err)
			}
			return nil
		})
		db.View(func(tx *bolt.Tx) error {
			buckets, err := NewProductBuckets(tx)
			if err != nil {
				t.Fatal(err)
			}
			got, err := buckets.List()
			if err != nil {
				t.Fatal(err)
			}

			want := products
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("got %v, want %v", got, want)
			}
			return nil
		})

	})
}

func TestProductBuckets_Update(t *testing.T) {

	t.Run("creates buckets if they don't exist", func(t *testing.T) {
		db := newTestDB(t)
		db.Update(func(tx *bolt.Tx) error {
			buckets, err := NewProductBuckets(tx)
			if err != nil {
				t.Fatal(err)
			}
			p := goshopify.Product{
				ID:     123,
				Handle: "product-123",
				Title:  "Product 123",
			}
			if err := buckets.Update(p); err != nil {
				t.Fatal(err)
			}
			if b := tx.Bucket([]byte("products.id")); b == nil {
				t.Fatal("expected products bucket to be created")
			}
			if b := tx.Bucket([]byte("products.title")); b == nil {
				t.Fatal("expected titles bucket to be created")
			}
			if b := tx.Bucket([]byte("products.handle")); b == nil {
				t.Fatal("expected handles bucket to be created")
			}
			return nil
		})
	})

	t.Run("allows product with ID only", func(t *testing.T) {
		db := newTestDB(t)
		db.Update(func(tx *bolt.Tx) error {
			buckets, err := NewProductBuckets(tx)
			if err != nil {
				t.Fatal(err)
			}
			p := goshopify.Product{ID: 123}
			if err := buckets.Update(p); err != nil {
				t.Fatal(err)
			}
			if b := tx.Bucket([]byte(bkeys.Products)); b == nil {
				t.Fatal("expected products bucket to be created")
			}
			if b := tx.Bucket([]byte(bkeys.ProductHandles)); b != nil {
				t.Fatal("expected no handles bucket")
			}
			if b := tx.Bucket([]byte(bkeys.ProductTitles)); b != nil {
				t.Fatal("expected no titles bucket")
			}
			return nil
		})
	})

	t.Run("can be called multiple times with same product", func(t *testing.T) {
		db := newTestDB(t)
		db.Update(func(tx *bolt.Tx) error {
			buckets, err := NewProductBuckets(tx)
			if err != nil {
				t.Fatal(err)
			}
			p := goshopify.Product{
				ID:     123,
				Handle: "product-123",
				Title:  "Product 123",
			}
			if err := buckets.Update(p); err != nil {
				t.Fatal(err)
			}
			if err := buckets.Update(p); err != nil {
				t.Fatal(err)
			}
			return nil
		})
	})

	t.Run("returns error if product handle is already set to another id", func(t *testing.T) {
		db := newTestDB(t)
		db.Update(func(tx *bolt.Tx) error {
			buckets, err := NewProductBuckets(tx)
			if err != nil {
				t.Fatal(err)
			}
			p1 := goshopify.Product{
				ID:     123,
				Handle: "product-123",
				Title:  "Product 123",
			}
			if err := buckets.Update(p1); err != nil {
				t.Fatal(err)
			}
			p2 := goshopify.Product{
				ID:     456,
				Handle: "product-123",
				Title:  "Product 123",
			}
			if err := buckets.Update(p2); err == nil {
				t.Fatalf("expected error, got %v", err)
			}
			return nil
		})
	})

	t.Run("returns error if product title is already set to another id", func(t *testing.T) {
		db := newTestDB(t)
		db.Update(func(tx *bolt.Tx) error {
			buckets, err := NewProductBuckets(tx)
			if err != nil {
				t.Fatal(err)
			}
			p1 := goshopify.Product{
				ID:     123,
				Handle: "product-123",
				Title:  "Product 456",
			}
			if err := buckets.Update(p1); err != nil {
				t.Fatal(err)
			}
			p2 := goshopify.Product{
				ID:     456,
				Handle: "product-456",
				Title:  "Product 456",
			}
			if err := buckets.Update(p2); err == nil {
				t.Fatalf("expected error, got %v", err)
			}
			return nil
		})
	})
}

func Test_setOnce(t *testing.T) {

	t.Run("sets initial value", func(t *testing.T) {
		k := []byte("foo")
		v := []byte("bar")

		updateTestBucket(t, func(b *bolt.Bucket) error {
			if err := setOnce(b, k, v); err != nil {
				return err
			}
			got := b.Get(k)
			if !bytes.Equal(got, v) {
				t.Fatalf("got %q, want %q", got, v)
			}
			return nil
		})
	})

	t.Run("does not return error when called multiple times with same value", func(t *testing.T) {
		k := []byte("foo")
		v := []byte("bar")

		updateTestBucket(t, func(b *bolt.Bucket) error {
			if err := setOnce(b, k, v); err != nil {
				t.Fatal(err)
			}
			if err := setOnce(b, k, v); err != nil {
				t.Fatal(err)
			}
			return nil
		})
	})

	t.Run("returns error when called multiple times with different value", func(t *testing.T) {
		k := []byte("foo")
		v1 := []byte("bar")
		v2 := []byte("baz")

		updateTestBucket(t, func(b *bolt.Bucket) error {
			if err := setOnce(b, k, v1); err != nil {
				t.Fatal(err)
			}
			if err := setOnce(b, k, v2); err == nil {
				t.Fatal("expected error but didn't get one")
			}
			return nil
		})
	})
}

func Test_int64ToBytes(t *testing.T) {
	tests := []struct {
		v    int64
		want []byte
	}{
		{
			v:    123456789,
			want: []byte("123456789"),
		},
		{
			v:    0,
			want: []byte("0"),
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.v), func(t *testing.T) {
			got := int64ToBytes(tt.v)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func updateTestBucket(t *testing.T, fn func(b *bolt.Bucket) error) error {
	db := newTestDB(t)

	return db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("test"))
		if err != nil {
			return err
		}
		return fn(bucket)
	})
}

func newTestDB(t *testing.T) *bolt.DB {
	filename := filepath.Join(t.TempDir(), "bolt.db")
	db, err := bolt.Open(filename, 0600, nil)
	if err != nil {
		t.Fatal(err)
		return nil
	}
	t.Cleanup(func() { db.Close() })
	return db
}
