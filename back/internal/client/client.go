package client

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/reinerRubin/froppyshima/back/internal"
	"github.com/reinerRubin/froppyshima/back/internal/engine"
	"github.com/reinerRubin/froppyshima/back/internal/protocol"
)

const eventBuffer = 100

func NewClient(conn *websocket.Conn, ct *ClientContext) *Client {
	gameEventsChan := make(chan *engine.GameEvent, eventBuffer)

	client := &Client{
		gameEvents: gameEventsChan,
		game: &ClientGame{
			Events:             gameEventsChan,
			gameRepository:     ct.GameRepository,
			playedGameRegister: ct.PlayedGameRegister,
		},

		connection: &clientConnection{
			conn: conn,
			out:  make(chan []byte, 256),
			in:   make(chan []byte),
		},
	}

	return client
}

type (
	ClientContext struct {
		GameRepository     internal.GameRepository
		PlayedGameRegister internal.PlayedGameRegister
	}

	Client struct {
		game       *ClientGame
		gameEvents chan *engine.GameEvent
		connection *clientConnection
	}
)

func (c *Client) Start() {
	go c.connection.readPump()
	go c.connection.writePump()
	go c.Run()
}

func (c *Client) Run() {
	defer c.game.Finalize()

	for {
		select {
		case rawCall, more := <-c.connection.in:
			if !more {
				return
			}

			call, err := protocol.ParseCall(rawCall)
			if err != nil {
				log.Printf("corrupt call: %s\n", err)
				continue
			}

			response := c.processCall(call)
			rawResponse, err := protocol.MarshallServerResponse(response)
			if err != nil {
				panic("we are doomed")
			}

			c.connection.out <- rawResponse
		case gameEvent, more := <-c.gameEvents:
			if !more {
				return
			}

			response, err := c.processEvent(gameEvent)
			if err != nil {
				log.Printf("corrupt event: %s", gameEvent)
			}

			rawResponse, err := protocol.MarshallServerResponse(response)
			if err != nil {
				panic("we are doomed")
			}

			c.connection.out <- rawResponse
		}
	}
}

func (c *Client) processEvent(e *engine.GameEvent) (*protocol.ServerResponse, error) {
	event := protocol.NewEvent()

	raw, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	event.Body = raw

	return event, nil
}

// TODO move to a register/call pattern (like HTTP) and add a separate marshal/unmoral layer
func (c *Client) processCall(call *protocol.Call) (response *protocol.ServerResponse) {
	response = protocol.NewCallResponse(call)

	switch call.Method {
	case protocol.MethodNew:
		gameID, err := c.game.New()
		if err != nil {
			response.Error = err.Error()
			return
		}

		resultResponse := protocol.NewResponse{
			ID: gameID,
		}
		rawResponse, err := json.Marshal(resultResponse)
		if err != nil {
			response.Error = err.Error()
			return
		}

		response.Body = rawResponse
	case protocol.MethodLoad:
		loadArgs := &protocol.LoadArgs{}
		err := json.Unmarshal(call.Body, loadArgs)
		if err != nil {
			response.Error = err.Error()
			return
		}

		info, err := c.game.Load(loadArgs.ID)
		if err != nil {
			response.Error = err.Error()
			return
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
			response.Error = err.Error()
			return
		}

		response.Body = rawResponse
	case protocol.MethodHit:
		hitArgs := &protocol.HitArgs{}
		err := json.Unmarshal(call.Body, hitArgs)
		if err != nil {
			response.Error = err.Error()
			return
		}

		result, err := c.game.Hit(&engine.Coord{
			X: hitArgs.X,
			Y: hitArgs.Y,
		})
		if err != nil {
			response.Error = err.Error()
			return
		}

		resultResponse := protocol.HitResponse{
			Result: result,
		}
		rawResponse, err := json.Marshal(resultResponse)
		if err != nil {
			response.Error = err.Error()
			return
		}

		response.Body = rawResponse
	default:
		response.Error = fmt.Errorf("unknown method: %s", call.Method).Error()
	}

	return
}
