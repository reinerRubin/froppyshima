package internal

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/reinerRubin/froppyshima/back/internal/engine"
)

const GameBucket = "game"

type GameRepository interface {
	Save(id engine.GameID, game *engine.Game) error
	Load(uuid engine.GameID) (*engine.Game, error)
}

type BoltDBGameRepository struct {
	dbProvider *BoltDBProvider
}

func NewBoltDBGameRepository(dbProvider *BoltDBProvider) (*BoltDBGameRepository, error) {
	gp := &BoltDBGameRepository{
		dbProvider: dbProvider,
	}

	err := gp.dbProvider.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(GameBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return gp, nil
}

func (gr *BoltDBGameRepository) Save(id engine.GameID, game *engine.Game) error {
	rawGame, err := json.Marshal(game)
	if err != nil {
		return err
	}

	return gr.dbProvider.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(GameBucket))
		return b.Put([]byte(id.String()), rawGame)
	})
}

func (gr *BoltDBGameRepository) Load(id engine.GameID) (*engine.Game, error) {
	var (
		rawGame []byte
		game    = new(engine.Game)
	)

	err := gr.dbProvider.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(GameBucket))
		rawGame = b.Get([]byte(id.String()))

		return nil
	})
	if err != nil {
		return nil, err
	}

	if len(rawGame) == 0 {
		return nil, fmt.Errorf(`game "%s" was not found`, id)
	}

	err = json.Unmarshal(rawGame, game)
	if err != nil {
		return nil, err
	}

	return game, nil
}
