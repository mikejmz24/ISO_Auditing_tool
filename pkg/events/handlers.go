package events

import (
	"ISO_Auditing_Tool/pkg/types"
	"context"
	"log"
)

func LoggingHandler() Handler {
	return func(ctx context.Context, event Event) error {
		log.Printf("Event: %s, Payload type: %T", event.Type, event.Payload)
		return nil
	}
}

// CreateMaterializedQueryEvent creates a MaterializedQueryCreated event from a MaterializedQuery
func CreateMaterializedQueryEvent(materializedQuery types.MaterializedQuery) Event {
	return NewMaterializedQueryCreatedEvent(
		materializedQuery.Name,
		materializedQuery.Definition,
		materializedQuery.Data, // This will be json.RawMessage but the function expects []byte
		materializedQuery.Version,
		materializedQuery.ErrorCount,
		materializedQuery.LastError,
	)
}

// RefreshMaterializedQueryEvent creates a MaterializedQueryRefreshRequested event from a MaterializedQuery
func RefreshMaterializedQueryEvent(materializedQuery types.MaterializedQuery) Event {
	return NewMaterializedQueryRefreshEvent(
		materializedQuery.Name,
		materializedQuery.Definition,
		materializedQuery.Data, // This will be json.RawMessage but the function expects []byte
		materializedQuery.Version,
		materializedQuery.ErrorCount,
		materializedQuery.LastError,
	)
}

// UpdateMaterializedQueryEvent creates a MaterializedQueryUpdated event from a MaterializedQuery
func UpdateMaterializedQueryEvent(materializedQuery types.MaterializedQuery) Event {
	return NewMaterializedQueryUpdatedEvent(
		materializedQuery.Name,
		materializedQuery.Definition,
		materializedQuery.Data, // This will be json.RawMessage but the function expects []byte
		materializedQuery.Version,
		materializedQuery.ErrorCount,
		materializedQuery.LastError,
	)
}

// TypedEventHandler is a wrapper that provides type-safe event handling
// It handles the type assertion for you to reduce boilerplate and errors
type TypedEventHandler struct {
	// Handler for data change events
	DataChangeHandler func(ctx context.Context, eventType EventType, payload DataChangePayload) error

	// Handler for materialized query events
	MaterializedQueryHandler func(ctx context.Context, eventType EventType, payload MaterializedQueryPayload) error

	// Handler for unknown event types
	FallbackHandler func(ctx context.Context, event Event) error
}

// HandleEvent implements the Handler function interface for TypedEventHandler
func (h TypedEventHandler) HandleEvent(ctx context.Context, event Event) error {
	switch event.Type {
	case DataCreated, DataUpdated, DataDeleted:
		if h.DataChangeHandler != nil {
			if payload, err := GetDataChangePayload(event); err == nil {
				return h.DataChangeHandler(ctx, event.Type, payload)
			} else {
				log.Printf("Error extracting DataChangePayload: %v", err)
			}
		}

	case MaterializedQueryCreated, MaterializedQueryUpdated, MaterializedQueryRefreshRequested:
		if h.MaterializedQueryHandler != nil {
			if payload, err := GetMaterializedQueryPayload(event); err == nil {
				return h.MaterializedQueryHandler(ctx, event.Type, payload)
			} else {
				log.Printf("Error extracting MaterializedQueryPayload: %v", err)
			}
		}
	}

	// Use fallback handler for unknown event types or if specific handler is not set
	if h.FallbackHandler != nil {
		return h.FallbackHandler(ctx, event)
	}

	// If we reach here, we didn't handle the event
	log.Printf("No handler for event type: %s", event.Type)
	return nil
}

// AsHandler converts a TypedEventHandler to a regular Handler function
func (h TypedEventHandler) AsHandler() Handler {
	return h.HandleEvent
}

// Example of creating a typed handler:
//
// func NewLoggingHandler() Handler {
//     return TypedEventHandler{
//         DataChangeHandler: func(ctx context.Context, eventType EventType, payload DataChangePayload) error {
//             log.Printf("Data change event: %s, Entity: %s, ID: %v",
//                 eventType, payload.EntityType, payload.EntityID)
//             return nil
//         },
//         MaterializedQueryHandler: func(ctx context.Context, eventType EventType, payload MaterializedQueryPayload) error {
//             log.Printf("Materialized query event: %s, Query: %s",
//                 eventType, payload.QueryName)
//             return nil
//         },
//     }.AsHandler()
// }
