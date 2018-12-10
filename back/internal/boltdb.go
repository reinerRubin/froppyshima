package internal

import (
	"os"

	"github.com/boltdb/bolt"
	"github.com/reinerRubin/froppyshima/back/internal/config"
)

const boltDBFileMode os.FileMode = 0600

type BoltDBProvider struct {
	DB *bolt.DB
}

func NewBoltDBProvider(config config.BoltDB) (*BoltDBProvider, error) {
	db, err := bolt.Open(config.Path, boltDBFileMode, nil)
	if err != nil {
		return nil, err
	}

	return &BoltDBProvider{
		DB: db,
	}, nil
}
