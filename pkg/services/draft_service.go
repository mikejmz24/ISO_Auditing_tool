// Contains draft business logic
// Calls the draft repository and applies transformations
package services

import (
	"ISO_Auditing_Tool/pkg/repositories"
	"ISO_Auditing_Tool/pkg/types"
	"context"
)

type DraftService struct {
	Repo repositories.DraftRepositoryInterface
}

// ensure DraftService implements DraftServiceInterface
var _ DraftServiceInterface = (*DraftService)(nil)

func NewDraftService(repo repositories.DraftRepositoryInterface) DraftServiceInterface {
	return &DraftService{Repo: repo}
}

func (s *DraftService) GetAll(ctx context.Context) ([]types.Draft, error) {
	return s.Repo.GetAllDrafts(ctx)
}

func (s *DraftService) Create(ctx context.Context, draft types.Draft) (types.Draft, error) {
	return s.Repo.CreateDraft(ctx, draft)
}

func (s *DraftService) Update(ctx context.Context, draft types.Draft) (types.Draft, error) {
	return s.Repo.UpdateDraft(ctx, draft)
}

func (s *DraftService) GetByID(ctx context.Context, draft types.Draft) (types.Draft, error) {
	return s.Repo.GetDraftByID(ctx, draft)
}

func (s *DraftService) Delete(ctx context.Context, draft types.Draft) (types.Draft, error) {
	return s.Repo.DeleteDraft(ctx, draft)
}

// List
