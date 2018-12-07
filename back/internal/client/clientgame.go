package client

import (
	"fmt"
	"log"

	"github.com/reinerRubin/froppyshima/back/internal"
	"github.com/reinerRubin/froppyshima/back/internal/engine"
)

type ClientGame struct {
	gameRepository     internal.GameRepository
	playedGameRegister internal.PlayedGameRegister

	Game   *engine.Game
	GameID engine.GameID

	Events chan *engine.GameEvent
}

func (cg *ClientGame) New() (engine.GameID, error) {
	var gameID engine.GameID

	cg.Game = engine.NewGame()
	err := cg.Game.Init()
	if err != nil {
		return gameID, err
	}

	if err := cg.Game.PutShips(); err != nil {
		return gameID, err
	}

	gameID = engine.NewGameID()
	if success := cg.RegisterPlay(gameID); !success {
		return gameID, fmt.Errorf("cant start game (already played): %s", gameID)
	}

	cg.GameID = gameID
	if err := cg.Save(); err != nil {
		cg.StopPlay(cg.GameID)
		return cg.GameID, err
	}

	return cg.GameID, nil
}

func (cg *ClientGame) Load(id engine.GameID) (*engine.GameMinInfo, error) {
	if success := cg.RegisterPlay(id); !success {
		return nil, fmt.Errorf("cant load game (already played): %s", id)
	}

	game, err := cg.gameRepository.Load(id)
	if err != nil {
		cg.StopPlay(id)
		return nil, err
	}

	cg.GameID = id
	cg.Game = game

	return cg.Game.MinInfo(), nil
}

func (cg *ClientGame) Hit(coord *engine.Coord) (result engine.HitResult, err error) {
	err = cg.Do(func() error {
		result, err = cg.Game.DotHit(coord)
		return err
	})

	return
}

func (cg *ClientGame) Do(fn func() error) error {
	defer func() {
		fmt.Println(cg.Game.Render())
	}()

	err := fn()
	if err != nil {
		return err
	}

	if err := cg.Save(); err != nil {
		return err
	}

	// TODO mv to separate process and make circular list
	events := cg.Game.PullEvents()
	for _, event := range events {
		select {
		case cg.Events <- event:
		default:
			log.Printf("event pull is full")
		}
	}

	return nil
}

func (cg *ClientGame) Save() error {
	return cg.gameRepository.Save(cg.GameID, cg.Game)
}

func (cg *ClientGame) RegisterPlay(id engine.GameID) bool {
	if success := cg.playedGameRegister.TryToPlay(id); !success {
		return false
	}

	return true
}

func (cg *ClientGame) StopPlay(id engine.GameID) {
	cg.playedGameRegister.StopPlay(id)
}

func (cg *ClientGame) Finalize() {
	cg.StopPlay(cg.GameID)
}
