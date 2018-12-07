package internal

import (
	"github.com/boltdb/bolt"
)

const DBPath = "/tmp/froppyshima.db"

type BoltDBProvider struct {
	DB *bolt.DB
}

func NewBoltDBProvider() (*BoltDBProvider, error) {
	db, err := bolt.Open(DBPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &BoltDBProvider{
		DB: db,
	}, nil
}
