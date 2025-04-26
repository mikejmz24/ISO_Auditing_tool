package events

import (
	"context"
	"fmt"
	"log"
	"sync"
)

// ErrorCallback defines a function that handles errors during event processing
type ErrorCallback func(eventType EventType, err error)

// DefaultErrorCallback provides a standard logging implementation for event errors
func DefaultErrorCallback(eventType EventType, err error) {
	log.Printf("Error handling event %s: %v", eventType, err)
}

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

// Publish synchronously publishes an event to all handlers
func (b *EventBus) Publish(ctx context.Context, event Event) error {
	b.mu.RLock()
	handlers, exists := b.handlers[event.Type]
	b.mu.RUnlock()

	if !exists {
		// No handlers registered for this event type
		return nil
	}

	var errs []error
	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			errs = append(errs, fmt.Errorf("handler error for event %s: %w", event.Type, err))
		}
	}

	if len(errs) > 0 {
		// Return just the first error for simplicity
		// In a more robust implementation, you might want to create a multi-error type
		return errs[0]
	}
	return nil
}

// AsyncPublish asynchronously publishes an event
// This maintains compatibility with existing code
func (b *EventBus) AsyncPublish(ctx context.Context, event Event) {
	b.AsyncPublishWithCallback(ctx, event, nil)
}

// AsyncPublishWithCallback asynchronously publishes an event with a custom error callback
func (b *EventBus) AsyncPublishWithCallback(ctx context.Context, event Event, errCallback ErrorCallback) {
	// Create a background context for continuing work after original context is done
	bgCtx := context.Background()

	go func() {
		if err := b.Publish(bgCtx, event); err != nil {
			if errCallback != nil {
				errCallback(event.Type, err)
			} else {
				DefaultErrorCallback(event.Type, err)
			}
		}
	}()
}

// AsyncPublishWithContext publishes an event asynchronously but respects the provided context
func (b *EventBus) AsyncPublishWithContext(ctx context.Context, event Event, errCallback ErrorCallback) {
	go func() {
		if err := b.Publish(ctx, event); err != nil {
			if errCallback != nil {
				errCallback(event.Type, err)
			} else {
				DefaultErrorCallback(event.Type, err)
			}
		}
	}()
}
