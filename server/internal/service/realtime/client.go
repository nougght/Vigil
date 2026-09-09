package realtime

import (
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"slices"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	realtime_model "github.com/nougght/monitoring-system/server/internal/model/realtime"
)

type Client struct {
	conn    *websocket.Conn
	toSend  chan *realtime_model.Message
	subs    map[string]struct{}
	subFunc func(subject string)
	mu      sync.RWMutex
}

func NewClient(conn *websocket.Conn, subFunc func(subject string)) *Client {
	return &Client{
		conn:    conn,
		toSend:  make(chan *realtime_model.Message, 256),
		subs:    make(map[string]struct{}),
		subFunc: subFunc,
	}
}

func (c *Client) Run() {
	c.runReader()
	c.runWriter()
}

func (c *Client) runReader() {
	go func() {
		defer func() {
			_ = c.conn.Close()
		}()

		for {
			var message realtime_model.ClientMessage
			err := c.conn.ReadJSON(&message)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
				}
				log.Printf("failed to read json ws: %s", err.Error())
				break
			}
			log.Printf("received client message: %#v", message)
			switch message.Type {
			case realtime_model.ClientMessageTypeAgentDetailed:
				msg := &realtime_model.AgentDetailedMessage{}
				err = json.Unmarshal(message.Payload, msg)
				if err != nil {
					log.Printf("failed to unmarshal client message ws: %s", err.Error())
					continue
				}
				log.Println("agent detailed message")
				c.mu.Lock()
				for _, agent := range msg.Agents {
					subj := fmt.Sprintf("%s.%s", realtime_model.ClientMessageTypeAgentDetailed, agent)
					c.subs[subj] = struct{}{}
					c.subFunc(subj)
				}
				c.mu.Unlock()

			}

		}
	}()
}

func (c *Client) runWriter() {
	go func() {
		ticker := time.NewTicker(time.Second * 2)
		defer func() {
			ticker.Stop()
			c.conn.Close()
		}()
		for {
			select {
			case message, ok := <-c.toSend:
				if !ok {
					// The hub closed the channel.
					c.conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}

				err := c.conn.WriteJSON(message)
				if err != nil {
					log.Printf("failed to write json message ws: %s", err.Error())
				}
			case <-ticker.C:
				if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}()
}

func (c *Client) GetSubs() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return slices.Collect(maps.Keys(c.subs))
}

func (c *Client) IsSub(subject string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.subs[subject]
	return ok
}
