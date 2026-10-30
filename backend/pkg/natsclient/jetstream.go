package natsclient

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/nougght/monitoring-system/gen/proto/messages"
	"google.golang.org/protobuf/proto"
)

type JetStreamClient struct {
	conn *nats.Conn
	js   jetstream.JetStream
}

func NewJetStreamClient(conn *nats.Conn, js jetstream.JetStream) *JetStreamClient {
	return &JetStreamClient{conn: conn, js: js}
}
func (c *JetStreamClient) InitStream() error {
	_, err := c.js.CreateOrUpdateStream(context.Background(), jetstream.StreamConfig{
		Name:      toStreamName(""),
		Subjects:  []string{},
		Storage:   jetstream.FileStorage,
		Retention: jetstream.LimitsPolicy,
		MaxAge:    24 * time.Hour,
		MaxMsgs:   100000,
	})
	return err
}

func (c *JetStreamClient) PublishRaw(ctx context.Context, subject string, data []byte) error {
	_, err := c.js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:      toStreamName(subject),
		Subjects:  []string{subject},
		Storage:   jetstream.FileStorage,
		Retention: jetstream.LimitsPolicy,
		MaxAge:    24 * time.Hour,
		MaxMsgs:   100000,
	})
	if err != nil {
		return fmt.Errorf("failed to create stream: %v", err)
	}

	_, err = c.js.Publish(ctx, subject, data)
	if err != nil {
		return fmt.Errorf("failed to publish: %v", err)
	}
	return nil
}

func (c *JetStreamClient) Publish(ctx context.Context, subject string, message *messages.EventMessage) error {
	data, err := proto.Marshal(message)
	if err != nil {
		log.Fatalf("failed to marshal proto: %v", err)
	}

	return c.PublishRaw(ctx, subject, data)
}

func (c *JetStreamClient) SubscribeRaw(ctx context.Context, subject string, handler func(ctx context.Context, data []byte) error) error {
	consumer, err := c.js.CreateOrUpdateConsumer(ctx, toStreamName(subject), jetstream.ConsumerConfig{
		Durable:       fmt.Sprintf("%s_%s", toStreamName(subject), uuid.NewString()), // Name that persists across restarts
		AckPolicy:     jetstream.AckExplicitPolicy,
		DeliverPolicy: jetstream.DeliverAllPolicy, // start from beginning
		AckWait:       30 * time.Minute,           // time to wait for ack
		MaxDeliver:    10,                         // max redelivery attempts
	})
	if err != nil {
		log.Fatalf("failed to create consumer: %v", err)
	}

	msgs, err := consumer.Messages()
	if err != nil {
		return fmt.Errorf("failed to start consumer: %v", err)
	}
	go func() {
		for {
			msg, err := msgs.Next(jetstream.NextContext(ctx))
			if err != nil {
				log.Println("failed to receive next jetstream message:", err)
				return
			}

			err = handler(ctx, msg.Data())
			if err != nil {
				_ = msg.Term()
				log.Println("failed to process jetstream message:", err)
				continue
			}

			_ = msg.Ack()
		}
	}()
	return nil
}

func (c *JetStreamClient) Subscribe(ctx context.Context, subject string, handler func(ctx context.Context, msg *messages.EventMessage) error) error {
	return c.SubscribeRaw(ctx, subject, func(ctx context.Context, data []byte) error {
		msg := &messages.EventMessage{}

		err := proto.Unmarshal(data, msg)
		if err != nil {
			log.Fatalf("failed to unmarshal proto: %v", err)
		}

		err = handler(ctx, msg)
		if err != nil {
			return err
		}

		return nil
	})
}

func toStreamName(subject string) string {
	return "events"

}
