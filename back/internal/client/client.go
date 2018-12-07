package client

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/reinerRubin/froppyshima/back/internal"
	"github.com/reinerRubin/froppyshima/back/internal/protocol"
)

const outMessageBuffer = 256

func NewClient(conn *websocket.Conn, ct *ClientContext) *Client {
	rawEvents := make(chan []byte, eventBuffer)

	client := &Client{
		rawEvents: rawEvents,
		clientWeb: NewClientWeb(ct, rawEvents),
		connection: &clientConnection{
			conn: conn,
			out:  make(chan []byte, outMessageBuffer),
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
		rawEvents chan []byte

		clientWeb  *ClientWeb
		connection *clientConnection
	}
)

// Start starts client related processes
// it's too naive approach and is not very scalable
func (c *Client) Start() {
	go c.connection.readPump()
	go c.connection.writePump()

	go c.clientWeb.Run()
	go c.Run()
}

func (c *Client) Run() {
	defer c.clientWeb.Stop()

	for {
		select {
		case rawCall, more := <-c.connection.in:
			if !more {
				return
			}

			call, err := protocol.ParseCall(rawCall)
			if err != nil {
				log.Printf("corrupt call: %s", err)
				continue
			}

			if err := c.processCall(call); err != nil {
				log.Printf("cant process call: %s", err)
			}
		case rawEvent, more := <-c.rawEvents:
			if !more {
				return
			}

			if err := c.processRawEvent(rawEvent); err != nil {
				log.Printf("cant process raw event: %s", err)
			}
		}
	}
}

func (c *Client) processRawEvent(raw []byte) error {
	event := protocol.NewEvent()
	event.Body = raw

	return c.SendResponse(event)
}

func (c *Client) processCall(call *protocol.Call) error {
	response := protocol.NewCallResponse(call)

	raw, err := NewClientWebRouter(c.clientWeb).Exec(call.Method, call.Body)
	if err != nil {
		response.Error = err.Error()
	} else {
		response.Body = raw
	}

	return c.SendResponse(response)
}

func (c *Client) SendResponse(response *protocol.ServerResponse) error {
	rawResponse, err := protocol.MarshallServerResponse(response)
	if err != nil {
		return err
	}

	select {
	case c.connection.out <- rawResponse:
	default:
		return fmt.Errorf("cant send to connection out: overflow")
	}

	return nil
}
