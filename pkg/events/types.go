package events

import (
	"context"
	"encoding/json"
	"fmt"
)

type EventType string
type EntityType string

const (
	EntityStandard    EntityType = "standard"
	EntityRequirement EntityType = "requirement"
	EntityQuestion    EntityType = "question"
	EntityEvidence    EntityType = "evidence"
)

type ChangeType string

const (
	ChangeCreated ChangeType = "created"
	ChangeUpdated ChangeType = "updated"
	ChangeDeleted ChangeType = "deleted"
)

const (
	EntityChanged                     EventType = "entity_changed"
	MaterializedQueryRefreshRequested EventType = "materialized_query_refresh_request"
	MaterializedQueryCreated          EventType = "materialized_query_created"
	MaterializedQueryUpdated          EventType = "materialized_query_updated"
)

const (
	DataCreated EventType = "data_created"
	DataUpdated EventType = "data_updated"
	DataDeleted EventType = "data_deleted"
)

type Event struct {
	Type    EventType
	Payload any
}

type Handler func(ctx context.Context, event Event) error

// Enhanced data change payload with stronger typing
type EntityChangePayload struct {
	EntityType    EntityType `json:"entity_type"`
	EntityID      any        `json:"entity_id"`
	ChangeType    ChangeType `json:"change_type"`
	ParentType    EntityType `json:"parent_type,omitempty"`
	ParentID      any        `json:"parent_id,omitempty"`
	AffectedQuery string     `json:"affected_query,omitempty"`
	Data          any        `json:"data,omitempty"` // Optional entity data for direct use
}

type MaterializedQueryPayload struct {
	QueryName       string          `json:"query_name"`
	QuerySQL        string          `json:"query_definition"`
	QueryData       json.RawMessage `json:"data"`
	QueryVersion    int             `json:"version"`
	QueryErrorCount int             `json:"error_count"`
	QueryLastError  string          `json:"last_error"`
}

// Create a unified function for entity change events
func NewEntityChangeEvent(
	entityType EntityType,
	entityID any,
	changeType ChangeType,
	affectedQuery string,
	parentType EntityType,
	parentID any,
	data any,
) Event {
	return Event{
		Type: EntityChanged,
		Payload: EntityChangePayload{
			EntityType:    entityType,
			EntityID:      entityID,
			ChangeType:    changeType,
			ParentType:    parentType,
			ParentID:      parentID,
			AffectedQuery: affectedQuery,
			Data:          data,
		},
	}
}

// Helper convenience functions for common entity types
func NewStandardEvent(standardID any, changeType ChangeType, affectedQuery string, data any) Event {
	return NewEntityChangeEvent(EntityStandard, standardID, changeType, affectedQuery, "", nil, data)
}

func NewRequirementEvent(requirementID any, changeType ChangeType, standardID any, affectedQuery string, data any) Event {
	return NewEntityChangeEvent(EntityRequirement, requirementID, changeType, affectedQuery, EntityStandard, standardID, data)
}

func NewQuestionEvent(questionID any, changeType ChangeType, requirementID any, affectedQuery string, data any) Event {
	return NewEntityChangeEvent(EntityQuestion, questionID, changeType, affectedQuery, EntityRequirement, requirementID, data)
}

func NewEvidenceEvent(evidenceID any, changeType ChangeType, questionID any, affectedQuery string, data any) Event {
	return NewEntityChangeEvent(EntityEvidence, evidenceID, changeType, affectedQuery, EntityQuestion, questionID, data)
}

// For backward compatibility
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

// NewMaterializedQueryCreatedEvent creates a MaterializedQueryCreated event
// Update the parameter type from []byte to json.RawMessage
func NewMaterializedQueryCreatedEvent(name string, sql string, data json.RawMessage, version int, errorCount int, lastError string) Event {
	return Event{
		Type: MaterializedQueryCreated,
		Payload: MaterializedQueryPayload{
			QueryName:       name,
			QuerySQL:        sql,
			QueryData:       data, // No conversion needed now
			QueryVersion:    version,
			QueryErrorCount: errorCount,
			QueryLastError:  lastError,
		},
	}
}

// GetEntityChangePayload extracts an EntityChangePayload from an event
func GetEntityChangePayload(event Event) (EntityChangePayload, error) {
	if event.Type != EntityChanged {
		return EntityChangePayload{}, fmt.Errorf("event type %s does not use EntityChangePayload", event.Type)
	}

	payload, ok := event.Payload.(EntityChangePayload)
	if !ok {
		return EntityChangePayload{}, fmt.Errorf("invalid payload type for event %s: expected EntityChangePayload, got %T",
			event.Type, event.Payload)
	}
	return payload, nil
}

// Backward compatibility functions
type DataChangePayload struct {
	EntityType    string
	EntityID      any
	ChangeType    string
	AffectedQuery string
}

func GetDataChangePayload(event Event) (DataChangePayload, error) {
	// For EntityChanged events, convert to the old format
	if event.Type == EntityChanged {
		entityPayload, ok := event.Payload.(EntityChangePayload)
		if ok {
			return DataChangePayload{
				EntityType:    string(entityPayload.EntityType),
				EntityID:      entityPayload.EntityID,
				ChangeType:    string(entityPayload.ChangeType),
				AffectedQuery: entityPayload.AffectedQuery,
			}, nil
		}
	}

	// Handle legacy events
	if event.Type == DataCreated || event.Type == DataUpdated || event.Type == DataDeleted {
		payload, ok := event.Payload.(DataChangePayload)
		if !ok {
			return DataChangePayload{}, fmt.Errorf("invalid payload type for event %s", event.Type)
		}
		return payload, nil
	}

	return DataChangePayload{}, fmt.Errorf("event type %s does not use DataChangePayload", event.Type)
}

type Subscriber interface {
	HandleEvent(ctx context.Context, event Event) (any, error)
}

// MaterializedQuerySource is an interface for objects that can be used to create materialized query events
type MaterializedQuerySource interface {
	GetName() string
	GetDefinition() string
	GetData() []byte // Keep as []byte to match current interface
	GetVersion() int
	GetErrorCount() int
	GetLastError() string
}

// Helper functions for creating specific events

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

// NewMaterializedQueryUpdatedEvent creates a MaterializedQueryUpdated event
// Update the parameter type from []byte to json.RawMessage
func NewMaterializedQueryUpdatedEvent(name string, sql string, data json.RawMessage, version int, errorCount int, lastError string) Event {
	return Event{
		Type: MaterializedQueryUpdated,
		Payload: MaterializedQueryPayload{
			QueryName:       name,
			QuerySQL:        sql,
			QueryData:       data, // No conversion needed now
			QueryVersion:    version,
			QueryErrorCount: errorCount,
			QueryLastError:  lastError,
		},
	}
}

// NewMaterializedQueryRefreshEvent creates a MaterializedQueryRefreshRequested event
// Update the parameter type from []byte to json.RawMessage
func NewMaterializedQueryRefreshEvent(name string, sql string, data json.RawMessage, version int, errorCount int, lastError string) Event {
	return Event{
		Type: MaterializedQueryRefreshRequested,
		Payload: MaterializedQueryPayload{
			QueryName:       name,
			QuerySQL:        sql,
			QueryData:       data, // No conversion needed now
			QueryVersion:    version,
			QueryErrorCount: errorCount,
			QueryLastError:  lastError,
		},
	}
}

// Helper functions to extract specific payload types

// GetDataChangePayload extracts a DataChangePayload from an event
//
//	func GetDataChangePayload(event Event) (DataChangePayload, error) {
//		// Check event type first to provide better error messages
//		switch event.Type {
//		case DataCreated, DataUpdated, DataDeleted:
//			// These event types should have DataChangePayload
//			break
//		default:
//			return DataChangePayload{}, fmt.Errorf("event type %s does not use DataChangePayload", event.Type)
//		}
//
//		payload, ok := event.Payload.(DataChangePayload)
//		if !ok {
//			return DataChangePayload{}, fmt.Errorf("invalid payload type for event %s: expected DataChangePayload, got %T",
//				event.Type, event.Payload)
//		}
//		return payload, nil
//	}
//
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
