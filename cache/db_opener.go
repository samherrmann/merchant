package cache

import (
	"path/filepath"

	bolt "go.etcd.io/bbolt"
)

type DBOpener interface {
	Open() (*bolt.DB, error)
}

func newDBOpener() (*dbOpener, error) {
	dir, err := directory()
	if err != nil {
		return nil, err
	}
	return &dbOpener{
		path: filepath.Join(dir, dbFilename),
	}, nil
}

type dbOpener struct {
	path string
}

func (p *dbOpener) Open() (*bolt.DB, error) {
	return bolt.Open(p.path, 0600, nil)
}
