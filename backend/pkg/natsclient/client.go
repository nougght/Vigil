package natsclient

import (
	"log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Client struct {
	//live    *nats.Conn
	durable *nats.Conn
	js      jetstream.JetStream
}

func NewClient(url string) *Client {
	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatalf("Failed to create JetStream context: %v", err)
	}
	return &Client{
		durable: nc,
		js:      js,
	}
}

//	func (c *Client) Live() *Live {
//		//return &Live{nc: c.live}
//	}
func (c *Client) Durable() *JetStreamClient {
	return &JetStreamClient{
		conn: c.durable,
		js:   c.js,
	}
}

//func (c *Client) Drain(ctx context.Context) error {
//	/* оба подключения */
//}
