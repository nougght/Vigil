package event

import (
	"context"
	"fmt"
	"log"

	"github.com/nougght/monitoring-system/gen/proto/messages"
)

const (
	EventSubjectPrefix            = "event"
	EventSystemEventSubjectPrefix = "event.system"

	AllEventsSubject = EventSubjectPrefix + ".>"
)

type EventProvider interface {
	Subscribe(ctx context.Context, subject string, handler func(ctx context.Context, msg *messages.EventMessage) error) error
}

type EventService struct {
	provider EventProvider
}

func NewEventService(provider EventProvider) *EventService {
	return &EventService{
		provider: provider,
	}
}

func (s *EventService) Start(ctx context.Context) error {
	log.Println("starting event service")
	err := s.provider.Subscribe(ctx, AllEventsSubject, func(ctx context.Context, event *messages.EventMessage) error {
		log.Printf("received event: %v", event)
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe all events: %w", err)
	}

	return nil
}
