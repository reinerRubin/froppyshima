package client

import (
	"encoding/json"
	"errors"
	"log"
	"sync"

	"github.com/reinerRubin/froppyshima/back/internal/engine"
	"github.com/reinerRubin/froppyshima/back/internal/protocol"
)

const eventBuffer = 100

func NewClientWebRouter(c *ClientWeb) MethodRouter {
	return MethodRouter{
		protocol.MethodNew:  c.NewGameHandler,
		protocol.MethodLoad: c.LoadGameHandler,
		protocol.MethodHit:  c.HitHandler,
	}
}

type ClientWeb struct {
	game *ClientGame

	gameEvents    <-chan *engine.GameEvent
	rawGameEvents chan<- []byte

	stopOnce       sync.Once
	stopChannel    chan struct{}
	stoppedChannel chan struct{}
}

func NewClientWeb(ct *ClientContext, rawEvents chan<- []byte) *ClientWeb {
	gameEventsChan := make(chan *engine.GameEvent, eventBuffer)

	client := &ClientWeb{
		rawGameEvents: rawEvents,
		gameEvents:    gameEventsChan,

		stopChannel:    make(chan struct{}),
		stoppedChannel: make(chan struct{}),

		game: &ClientGame{
			Events:             gameEventsChan,
			gameRepository:     ct.GameRepository,
			playedGameRegister: ct.PlayedGameRegister,
		},
	}

	return client
}

func (c *ClientWeb) NewGameHandler(params json.RawMessage) (json.RawMessage, error) {
	gameID, err := c.game.New()
	if err != nil {
		return nil, err
	}

	resultResponse := protocol.NewResponse{
		ID: gameID,
	}
	rawResponse, err := json.Marshal(resultResponse)
	if err != nil {
		return nil, err
	}

	return rawResponse, nil
}

func (c *ClientWeb) LoadGameHandler(params json.RawMessage) (json.RawMessage, error) {
	loadArgs := &protocol.LoadArgs{}
	err := json.Unmarshal(params, loadArgs)
	if err != nil {
		return nil, err
	}

	info, err := c.game.Load(loadArgs.ID)
	if err != nil {
		return nil, err
	}

	loadResponse := &protocol.LoadResponse{
		Maxx: info.Maxx,
		Maxy: info.Maxy,

		FailHits:    info.FailHits,
		SuccessHits: info.SuccessHits,
		FatalHits:   info.FatalHits,

		AnyMoreShips: info.AnyMoreShips,
	}

	rawResponse, err := json.Marshal(loadResponse)
	if err != nil {
		return nil, err
	}

	return rawResponse, nil
}

func (c *ClientWeb) HitHandler(params json.RawMessage) (json.RawMessage, error) {
	hitArgs := &protocol.HitArgs{}
	err := json.Unmarshal(params, hitArgs)
	if err != nil {
		return nil, err
	}

	result, err := c.game.Hit(&engine.Coord{
		X: hitArgs.X,
		Y: hitArgs.Y,
	})
	if err != nil {
		return nil, err
	}

	resultResponse := protocol.HitResponse{
		Result: result,
	}
	rawResponse, err := json.Marshal(resultResponse)
	if err != nil {
		return nil, err
	}

	return rawResponse, nil
}

func (c *ClientWeb) Stop() {
	c.stopOnce.Do(func() {
		close(c.stopChannel)
	})

	<-c.stoppedChannel
}

func (c *ClientWeb) Run() {
	defer func() {
		close(c.stoppedChannel)
		close(c.rawGameEvents)
		c.game.Finalize()
	}()

	for {
		select {
		case gameEvent, more := <-c.gameEvents:
			if !more {
				return
			}

			err := c.processEvent(gameEvent)
			if err != nil {
				log.Printf("cant process event: %s", gameEvent)
			}
		case <-c.stopChannel:
			return
		}
	}
}

func (c *ClientWeb) processEvent(e *engine.GameEvent) error {
	raw, err := json.Marshal(e)
	if err != nil {
		return err
	}

	select {
	case c.rawGameEvents <- raw:
	default:
		return errors.New("rawGameEvents is overflowed")
	}

	return nil
}
