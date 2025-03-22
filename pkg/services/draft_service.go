// Contains draft business logic
// Calls the draft repository and applies transformations
package services

import (
	"ISO_Auditing_Tool/pkg/repositories"
	"ISO_Auditing_Tool/pkg/types"
	"context"
)

type DraftService struct {
	Repo repositories.DraftRepository
}

func NewDraftService(repo repositories.DraftRepository) *DraftService {
	return &DraftService{Repo: repo}
}

func (s *DraftService) Create(ctx context.Context, draft types.Draft) (types.Draft, error) {
	return s.Repo.CreateDraft(ctx, draft)

	// Create
	// GetByID
	// Update
	// Delete
	// List
}

func (s *DraftService) Update(ctx context.Context, draft types.Draft) (types.Draft, error) {
	return s.Repo.UpdateDraft(ctx, draft)
}
