package client

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type clientConnection struct {
	conn *websocket.Conn
	out  chan []byte
	in   chan []byte
}

var expectedCodes = []int{websocket.CloseGoingAway, websocket.CloseAbnormalClosure}

func (c *clientConnection) read() {
	defer func() {
		c.conn.Close()
		close(c.in)
	}()

	{ // init read
		if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			return
		}
		c.conn.SetReadLimit(maxMessageSize)
		c.conn.SetPongHandler(func(string) error {
			return c.conn.SetReadDeadline(time.Now().Add(pongWait))
		})
	}

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, expectedCodes...) {
				log.Printf("error: %v", err)
			}
			break
		}

		c.in <- message
	}
}
func (c *clientConnection) write() {
	if err := c.runWrite(); err != nil {
		log.Printf("write error: %s", err)
	}
}

func (c *clientConnection) runWrite() error {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.out:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return err
			}
			if !ok {
				err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return err
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return err
			}
			if _, err = w.Write(message); err != nil {
				return err
			}

			if err := w.Close(); err != nil {
				return err
			}
		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return err
			}

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return err
			}
		}
	}
}
