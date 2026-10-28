package natsclient

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Client struct {
	live    *nats.Conn
	durable *nats.Conn
	js      jetstream.JetStream
}

//func (c *Client) Live() *Live {
//	//return &Live{nc: c.live}
//}
//func (c *Client) Durable() *Durable {
//	//return &Durable{js: c.js}
//}
//func (c *Client) Drain(ctx context.Context) error {
//	/* оба подключения */
//}

type Live struct {
}
type Durable struct {
}
