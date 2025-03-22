package events

import (
	"context"
	"fmt"
	"sync"
)

type EventBus struct {
	mu       sync.RWMutex
	handlers map[EventType][]Handler
}

func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[EventType][]Handler),
	}
}

func (b *EventBus) Subscribe(eventType EventType, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.handlers == nil {
		b.handlers = make(map[EventType][]Handler)
	}

	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

func (b *EventBus) SubscribeAll(handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for eventType := range b.handlers {
		b.handlers[eventType] = append(b.handlers[eventType], handler)
	}
}

func (b *EventBus) Publish(ctx context.Context, event Event) error {
	b.mu.RLock()
	handlers := b.handlers[event.Type]
	b.mu.RUnlock()

	var lastErr error
	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			lastErr = err
			// Consider logging the error but continuing to process other handlers
		}
	}

	return lastErr
}

func (b *EventBus) AsyncPublish(ctx context.Context, event Event) {
	go func() {
		if err := b.Publish(ctx, event); err != nil {
			// Log the erro but don't propagete it sinve this is AsyncPublish
			fmt.Printf("Error publishing event %s: %v\n", event.Type, err)
		}
	}()
}
