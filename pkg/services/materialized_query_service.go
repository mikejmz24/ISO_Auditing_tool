// Contains materialized query business logic
// Calls the materalized query repository and applies transformations
package services

import (
	"ISO_Auditing_Tool/pkg/custom_errors"
	"ISO_Auditing_Tool/pkg/events"
	"ISO_Auditing_Tool/pkg/repositories"
	"ISO_Auditing_Tool/pkg/types"
	"context"
	"errors"
	"fmt"
)

type MaterializedQueryService struct {
	Repo     repositories.MaterializedQueryRepository
	eventBus *events.EventBus
}

func NewMaterializedQueryService(repo repositories.MaterializedQueryRepository, eventBus *events.EventBus) *MaterializedQueryService {
	service := &MaterializedQueryService{
		Repo:     repo,
		eventBus: eventBus}

	eventBus.Subscribe(events.MaterializedQueryCreated, service.handleMaterializedQueryCreated)
	eventBus.Subscribe(events.MaterializedQueryRefreshRequested, service.handleRefreshRequest)

	return service
}

func (s *MaterializedQueryService) GetByName(ctx context.Context, name string) (types.MaterializedQuery, error) {
	return s.Repo.GetByNameMaterializedQuery(ctx, name)
}

func (s *MaterializedQueryService) handleMaterializedQueryCreated(ctx context.Context, event events.Event) error {
	payload, err := events.GetMaterializedQueryPayload(event)
	if err != nil {
		return err
	}

	// Create a new query with the payload data
	query := types.MaterializedQuery{
		Name:       payload.QueryName,
		Definition: payload.QuerySQL,
		Data:       payload.QueryData,
		Version:    payload.QueryVersion,
		ErrorCount: payload.QueryErrorCount,
		LastError:  payload.QueryLastError,
	}

	// Check if the query already exists
	existingQuery, err := s.Repo.GetByNameMaterializedQuery(ctx, payload.QueryName)

	if err != nil {
		// If it's a not found error, create a new query
		if errors.Is(err, custom_errors.ErrNotFound) {
			_, err = s.Repo.CreateMaterializedQuery(ctx, query)
			if err != nil {
				return fmt.Errorf("failed to create materialized query: %w", err)
			}
		} else {
			// For any other error, return it
			return fmt.Errorf("error checking for existing materialized query: %w", err)
		}
	} else {
		// Query exists, update it with new values but preserve the ID
		query.ID = existingQuery.ID               // Preserve the ID of the existing query
		query.Version = existingQuery.Version + 1 // Increment version

		_, err = s.Repo.UpdateMaterializedQuery(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to update materialized query: %w", err)
		}
	}

	// Refresh the query after creating or updating
	return s.RefreshMaterializedQuery(ctx, payload.QueryName)
}

func (s *MaterializedQueryService) handleRefreshRequest(ctx context.Context, event events.Event) error {
	payload, ok := event.Payload.(events.MaterializedQueryPayload)
	if !ok {
		return fmt.Errorf("Invalid payload type for RefreshRequested event")
	}

	return s.RefreshMaterializedQuery(ctx, payload.QueryName)
}

func (s *MaterializedQueryService) RefreshMaterializedQuery(ctx context.Context, name string) error {
	query, err := s.Repo.GetByNameMaterializedQuery(ctx, name)
	if err != nil {
		return fmt.Errorf("Failed to get materialized query: %w", err)
	}

	_, err = s.Repo.UpdateMaterializedQuery(ctx, query)
	return err
}

func (s *MaterializedQueryService) PublishEvent(ctx context.Context, event events.Event) error {
	return s.eventBus.Publish(ctx, event)
}
