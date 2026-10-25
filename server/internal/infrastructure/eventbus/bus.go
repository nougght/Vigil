package eventbus

import (
	"context"
	"errors"
	"log"
	"slices"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/nougght/monitoring-system/server/internal/config"
	model "github.com/nougght/monitoring-system/server/internal/model/event"
)

type SubID uuid.UUID

type SubChans map[SubID]chan<- model.Event

var (
	ErrInvalidSubjecName       = errors.New("invalid subject name")
	ErrFullSubjectNameRequired = errors.New("full subject name required")
)

type EventBus struct {
	cfg  *config.Config
	mu   sync.RWMutex
	subs map[string]SubChans

	busCtx       context.Context
	cancelBusCtx context.CancelFunc
	once         sync.Once
	wg           sync.WaitGroup
}

func NewEventBus(cfg *config.Config) *EventBus {
	return &EventBus{
		cfg:  cfg,
		subs: make(map[string]SubChans),
	}
}

func (b *EventBus) Start(ctx context.Context) {
	b.busCtx, b.cancelBusCtx = context.WithCancel(ctx)
}

func (b *EventBus) Shutdown(ctx context.Context) error {
	if b.cancelBusCtx == nil {
		return errors.New("bus not started")
	}
	b.once.Do(func() {
		b.cancelBusCtx()
		b.mu.Lock()
		defer b.mu.Unlock()
		for _, subs := range b.subs {
			for _, s := range subs {
				close(s)
			}
		}
		b.subs = make(map[string]SubChans)
	})

	wgDone := make(chan struct{})

	go func() {
		b.wg.Wait()
		close(wgDone)
	}()

	var err error
	select {
	case <-wgDone:
	case <-ctx.Done():
		err = ctx.Err()
	}

	return err
}

func (b *EventBus) Subscribe(subject string, handler model.EventHandler, bufferSize int, droppedCount *int) (uuid.UUID, error) {
	if b.cancelBusCtx == nil {
		return uuid.Nil, errors.New("bus not started")
	}
	log.Printf("sub event: %s", subject)
	if !isValidSubject(subject) {
		return uuid.Nil, ErrInvalidSubjecName
	}
	if bufferSize == 0 {
		bufferSize = 1
	}
	ch := make(chan model.Event, bufferSize)

	id := uuid.New()
	b.mu.Lock()
	subs, ok := b.subs[subject]
	if !ok { // new map for new subject
		subs = make(SubChans)
		b.subs[subject] = subs
	}
	subs[SubID(id)] = ch
	b.mu.Unlock()

	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		for {
			select {
			case <-b.busCtx.Done():
				return

			case event, ok := <-ch:
				if !ok {
					return
				}
				handler(b.busCtx, event)
			}
		}
	}()
	return id, nil
}

func (b *EventBus) Unsubscribe(subject string, id uuid.UUID) error {
	log.Printf("unsub event: %s", subject)
	if !isValidSubject(subject) {
		return ErrInvalidSubjecName
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	subs, ok := b.subs[subject]
	if !ok {
		log.Printf("no such subject to unsub: %s", subject)
		return nil
	}
	ch, ok := subs[SubID(id)]
	if !ok {
		log.Printf("no such sub to unsub: %s", subject)
		return nil
	}
	close(ch)
	delete(subs, SubID(id))
	if len(subs) == 0 {
		delete(b.subs, subject)
	}
	return nil
}

func (b *EventBus) Publish(ctx context.Context, event model.Event) error {
	log.Printf("pub event: %s", event.Subject())
	if !isFullSubject(event.Subject()) {
		return ErrFullSubjectNameRequired
	}

	//for subject, ch := range b.subs {
	//
	//}
	//isWildcard, tokens := isWildcard(event.Subject())
	//if !isWildcard {
	b.mu.RLock()
	defer b.mu.RUnlock()
	chans, ok := b.subs[event.Subject()]
	if ok {
		for _, ch := range chans {
			select {
			case ch <- event:
			default:
			}
		}
	}
	target := strings.Split(event.Subject(), ".")
	for subject, chans := range b.subs {
		ok, tokens := isWildcard(subject)
		if ok && matchSubject(tokens, target) {
			for _, ch := range chans {
				select {
				case ch <- event:
				default:
				}
			}
		}
	}
	return nil
}

func isValidSubject(subject string) bool {
	tokens := strings.Split(subject, ".")
	return isValidSubjectByTokens(tokens)
}

func isValidSubjectByTokens(subject []string) bool {
	for i, token := range subject {
		switch {
		case token == "":
			return false
		case token == ">": // only at end
			if i != len(subject)-1 {
				return false
			}
		case token == "*":
		case strings.ContainsAny(token, ">* \t\r\n"):
			// can't contain wildcard or space inside token
			return false
		}
	}
	return true
}

func isFullSubject(subject string) bool {
	return !strings.ContainsAny(subject, ">*")
}

func isWildcard(subject string) (bool, []string) {
	tokens := strings.Split(subject, ".")
	if tokens[len(tokens)-1] == ">" {
		return true, tokens
	}
	if slices.Contains(tokens, "*") {
		return true, tokens
	}
	return false, nil
}

func matchSubject(subject, target []string) bool {
	for i, token := range subject {
		// can't be longer than target
		if i >= len(target) {
			return false
		}
		// skip all remaining tokens
		if token == ">" {
			return true
		}
		// skip current token
		if token != "*" && token != target[i] {
			return false
		}
	}
	// if all subject tokens match, check is there remaining target tokens
	return len(subject) == len(target)
}
