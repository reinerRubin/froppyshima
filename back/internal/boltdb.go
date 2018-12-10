package internal

import (
	"github.com/boltdb/bolt"
	"github.com/reinerRubin/froppyshima/back/internal/config"
)

const DBPath = "/tmp/froppyshima.db"

type BoltDBProvider struct {
	DB *bolt.DB
}

func NewBoltDBProvider(config config.BoltDB) (*BoltDBProvider, error) {
	db, err := bolt.Open(config.Path, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &BoltDBProvider{
		DB: db,
	}, nil
}
