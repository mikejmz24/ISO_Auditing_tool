// pkg/events/handlers.go
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

// EntityChangeHandler creates a handler for entity change events
func EntityChangeHandler(entityService EntityService) Handler {
	return func(ctx context.Context, event Event) error {
		// Handle both new and legacy event types
		if event.Type == EntityChanged {
			payload, err := GetEntityChangePayload(event)
			if err != nil {
				log.Printf("Error extracting payload from entity change event: %v", err)
				return err
			}

			log.Printf("Processing entity change: %s ID=%v, ChangeType=%s",
				payload.EntityType, payload.EntityID, payload.ChangeType)

			return entityService.HandleEntityChange(ctx, payload)
		} else if event.Type == DataCreated || event.Type == DataUpdated || event.Type == DataDeleted {
			// Legacy event types
			legacyPayload, err := GetDataChangePayload(event)
			if err != nil {
				log.Printf("Error extracting payload from legacy event: %v", err)
				return err
			}

			// Convert to new format
			payload := EntityChangePayload{
				EntityType:    EntityType(legacyPayload.EntityType),
				EntityID:      legacyPayload.EntityID,
				ChangeType:    ChangeType(legacyPayload.ChangeType),
				AffectedQuery: legacyPayload.AffectedQuery,
			}

			log.Printf("Processing legacy entity change: %s ID=%v, ChangeType=%s",
				payload.EntityType, payload.EntityID, payload.ChangeType)

			return entityService.HandleEntityChange(ctx, payload)
		}

		return nil
	}
}

// MaterializedQueryHandler creates a handler for materialized query events
func MaterializedQueryHandler(materializedQueryService MaterializedQueryService) Handler {
	return func(ctx context.Context, event Event) error {
		payload, err := GetMaterializedQueryPayload(event)
		if err != nil {
			log.Printf("Error extracting payload from materialized query event: %v", err)
			return err
		}

		log.Printf("Processing materialized query event: %s for query %s",
			event.Type, payload.QueryName)

		return materializedQueryService.HandleQueryEvent(ctx, event.Type, payload)
	}
}

// HTMLCacheHandler creates a handler for updating HTML cache based on materialized query updates
func HTMLCacheHandler(htmlCacheService HTMLCacheService) Handler {
	return func(ctx context.Context, event Event) error {
		payload, err := GetMaterializedQueryPayload(event)
		if err != nil {
			log.Printf("Error extracting payload from materialized query event: %v", err)
			return err
		}

		log.Printf("Processing HTML cache update for query: %s", payload.QueryName)

		return htmlCacheService.RefreshHTMLForQuery(ctx, payload.QueryName)
	}
}

// CreateMaterializedQueryEvent creates a MaterializedQueryCreated event from a MaterializedQuery
func CreateMaterializedQueryEvent(materializedQuery types.MaterializedJSONQuery) Event {
	return NewMaterializedQueryCreatedEvent(
		materializedQuery.Name,
		"", // Definition isn't needed here
		materializedQuery.Data,
		materializedQuery.Version,
		materializedQuery.ErrorCount,
		materializedQuery.LastError,
	)
}

// UpdateMaterializedQueryEvent creates a MaterializedQueryUpdated event from a MaterializedQuery
func UpdateMaterializedQueryEvent(materializedQuery types.MaterializedJSONQuery) Event {
	return NewMaterializedQueryUpdatedEvent(
		materializedQuery.Name,
		"", // Definition isn't needed here
		materializedQuery.Data,
		materializedQuery.Version,
		materializedQuery.ErrorCount,
		materializedQuery.LastError,
	)
}

// Service interfaces that work with these handlers

// EntityService handles entity changes with the optimized payload format
type EntityService interface {
	HandleEntityChange(ctx context.Context, payload EntityChangePayload) error
}

// MaterializedQueryService handles materialized query operations
type MaterializedQueryService interface {
	HandleQueryEvent(ctx context.Context, eventType EventType, payload MaterializedQueryPayload) error
}

// HTMLCacheService handles HTML cache operations
type HTMLCacheService interface {
	RefreshHTMLForQuery(ctx context.Context, queryName string) error
}

// TypedEventHandler for handling typed events (keeping this from your original file, but updated)
type TypedEventHandler struct {
	// Handler for entity change events (new format)
	EntityChangeHandler func(ctx context.Context, eventType EventType, payload EntityChangePayload) error

	// Handler for legacy data change events
	DataChangeHandler func(ctx context.Context, eventType EventType, payload DataChangePayload) error

	// Handler for materialized query events
	MaterializedQueryHandler func(ctx context.Context, eventType EventType, payload MaterializedQueryPayload) error

	// Handler for unknown event types
	FallbackHandler func(ctx context.Context, event Event) error
}

// HandleEvent implements the Handler function interface for TypedEventHandler
func (h TypedEventHandler) HandleEvent(ctx context.Context, event Event) error {
	switch event.Type {
	case EntityChanged:
		if h.EntityChangeHandler != nil {
			if payload, err := GetEntityChangePayload(event); err == nil {
				return h.EntityChangeHandler(ctx, event.Type, payload)
			} else {
				log.Printf("Error extracting EntityChangePayload: %v", err)
			}
		}

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

// Example usage: Create a MaterializedJSONService that implements EntityService
type MaterializedJSONServiceExample struct {
	// Repository fields would go here
}

func (s *MaterializedJSONServiceExample) HandleEntityChange(ctx context.Context, payload EntityChangePayload) error {
	// Example implementation
	switch payload.EntityType {
	case EntityStandard:
		// Handle standard changes
		standardID, ok := payload.EntityID.(int)
		if !ok {
			return nil
		}

		switch payload.ChangeType {
		case ChangeCreated, ChangeUpdated:
			// Update standard materialized query
			return s.updateStandardMaterializedQuery(ctx, standardID)
		case ChangeDeleted:
			// Delete standard materialized query
			return s.deleteStandardMaterializedQuery(ctx, standardID)
		}

	case EntityRequirement:
		// Handle requirement changes
		requirementID, ok := payload.EntityID.(int)
		if !ok {
			return nil
		}

		// Use the parent ID to determine which standard needs updating
		standardID, ok := payload.ParentID.(int)
		if !ok {
			return nil
		}

		// Update related materialized queries
		return s.updateRequirementMaterializedQuery(ctx, requirementID, standardID)

		// Similar cases for questions and evidence
	}

	return nil
}

// These methods would be implemented by the real service
func (s *MaterializedJSONServiceExample) updateStandardMaterializedQuery(ctx context.Context, standardID int) error {
	return nil
}

func (s *MaterializedJSONServiceExample) deleteStandardMaterializedQuery(ctx context.Context, standardID int) error {
	return nil
}

func (s *MaterializedJSONServiceExample) updateRequirementMaterializedQuery(ctx context.Context, requirementID, standardID int) error {
	return nil
}
