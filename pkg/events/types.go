package events

import (
	"context"
	"fmt"
)

type EventType string

const (
	DataCreated EventType = "data_created"
	DataUpdated EventType = "data_updated"
	DataDeleted EventType = "data_deleted"

	MaterializedQueryRefreshRequested EventType = "materialized_query_refresh_request"
	MaterializedQueryCreated          EventType = "materialized_query_created"
	MaterializedQueryUpdated          EventType = "materialized_query_updated"
)

type Event struct {
	Type    EventType
	Payload any
}

type Handler func(ctx context.Context, event Event) error

type Subscriber interface {
	HandleEvent(ctx context.Context, event Event) (any, error)
}

type DataChangePayload struct {
	EntityType    string
	EntityID      any
	ChangeType    string
	AffectedQuery string
}

type MaterializedQueryPayload struct {
	QueryName       string `json:"query_name"`
	QuerySQL        string `json:"query_definition"`
	QueryData       []byte `json:"data"`
	QueryVersion    int    `json:"version"`
	QueryErrorCount int    `json:"error_count"`
	QueryLastError  string `json:"last_error"`
}

// Helper functions for creating specific events

// NewDataCreatedEvent creates a DataCreated event
func NewDataCreatedEvent(entityType string, entityID any, affectedQuery string) Event {
	return Event{
		Type: DataCreated,
		Payload: DataChangePayload{
			EntityType:    entityType,
			EntityID:      entityID,
			ChangeType:    "created",
			AffectedQuery: affectedQuery,
		},
	}
}

// NewDataUpdatedEvent creates a DataUpdated event
func NewDataUpdatedEvent(entityType string, entityID any, affectedQuery string) Event {
	return Event{
		Type: DataUpdated,
		Payload: DataChangePayload{
			EntityType:    entityType,
			EntityID:      entityID,
			ChangeType:    "updated",
			AffectedQuery: affectedQuery,
		},
	}
}

// NewDataDeletedEvent creates a DataDeleted event
func NewDataDeletedEvent(entityType string, entityID any, affectedQuery string) Event {
	return Event{
		Type: DataDeleted,
		Payload: DataChangePayload{
			EntityType:    entityType,
			EntityID:      entityID,
			ChangeType:    "deleted",
			AffectedQuery: affectedQuery,
		},
	}
}

// NewMaterializedQueryCreatedEvent creates a MaterializedQueryCreated event
func NewMaterializedQueryCreatedEvent(name string, sql string, data []byte, version int, errorCount int, lastError string) Event {
	return Event{
		Type: MaterializedQueryCreated,
		Payload: MaterializedQueryPayload{
			QueryName:       name,
			QuerySQL:        sql,
			QueryData:       data,
			QueryVersion:    version,
			QueryErrorCount: errorCount,
			QueryLastError:  lastError,
		},
	}
}

// NewMaterializedQueryUpdatedEvent creates a MaterializedQueryUpdated event
func NewMaterializedQueryUpdatedEvent(name string, sql string, data []byte, version int, errorCount int, lastError string) Event {
	return Event{
		Type: MaterializedQueryUpdated,
		Payload: MaterializedQueryPayload{
			QueryName:       name,
			QuerySQL:        sql,
			QueryData:       data,
			QueryVersion:    version,
			QueryErrorCount: errorCount,
			QueryLastError:  lastError,
		},
	}
}

// NewMaterializedQueryRefreshEvent creates a MaterializedQueryRefreshRequested event
func NewMaterializedQueryRefreshEvent(name string, sql string, data []byte, version int, errorCount int, lastError string) Event {
	return Event{
		Type: MaterializedQueryRefreshRequested,
		Payload: MaterializedQueryPayload{
			QueryName:       name,
			QuerySQL:        sql,
			QueryData:       data,
			QueryVersion:    version,
			QueryErrorCount: errorCount,
			QueryLastError:  lastError,
		},
	}
}

// Helper functions to extract specific payload types

// GetDataChangePayload extracts a DataChangePayload from an event
func GetDataChangePayload(event Event) (DataChangePayload, error) {
	// Check event type first to provide better error messages
	switch event.Type {
	case DataCreated, DataUpdated, DataDeleted:
		// These event types should have DataChangePayload
		break
	default:
		return DataChangePayload{}, fmt.Errorf("event type %s does not use DataChangePayload", event.Type)
	}

	payload, ok := event.Payload.(DataChangePayload)
	if !ok {
		return DataChangePayload{}, fmt.Errorf("invalid payload type for event %s: expected DataChangePayload, got %T",
			event.Type, event.Payload)
	}
	return payload, nil
}

// GetMaterializedQueryPayload extracts a MaterializedQueryPayload from an event
func GetMaterializedQueryPayload(event Event) (MaterializedQueryPayload, error) {
	// Check event type first to provide better error messages
	switch event.Type {
	case MaterializedQueryCreated, MaterializedQueryUpdated, MaterializedQueryRefreshRequested:
		// These event types should have MaterializedQueryPayload
		break
	default:
		return MaterializedQueryPayload{}, fmt.Errorf("event type %s does not use MaterializedQueryPayload", event.Type)
	}

	payload, ok := event.Payload.(MaterializedQueryPayload)
	if !ok {
		return MaterializedQueryPayload{}, fmt.Errorf("invalid payload type for event %s: expected MaterializedQueryPayload, got %T",
			event.Type, event.Payload)
	}
	return payload, nil
}

// ValidateEventPayload checks if an event has the correct payload type
// This can be used in tests or when receiving events from external sources
func ValidateEventPayload(event Event) error {
	switch event.Type {
	case DataCreated, DataUpdated, DataDeleted:
		_, ok := event.Payload.(DataChangePayload)
		if !ok {
			return fmt.Errorf("invalid payload type for event %s: expected DataChangePayload, got %T",
				event.Type, event.Payload)
		}
	case MaterializedQueryCreated, MaterializedQueryUpdated, MaterializedQueryRefreshRequested:
		_, ok := event.Payload.(MaterializedQueryPayload)
		if !ok {
			return fmt.Errorf("invalid payload type for event %s: expected MaterializedQueryPayload, got %T",
				event.Type, event.Payload)
		}
	default:
		return fmt.Errorf("unknown event type: %s", event.Type)
	}
	return nil
}
