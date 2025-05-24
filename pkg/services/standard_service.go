// Contains standard business logic for rendering HTML
// Calls the standard repository and applies transformations
package services

import (
	"ISO_Auditing_Tool/pkg/repositories"
	"ISO_Auditing_Tool/pkg/types"
	"context"
)

type StandardService struct {
	Repo repositories.StandardRepositoryInterface
}

func NewStandardService(repo repositories.StandardRepositoryInterface) *StandardService {
	return &StandardService{Repo: repo}
}

func (s *StandardService) GetByID(ctx context.Context, standard types.Standard) (types.Standard, error) {
	return s.Repo.GetByIDStandard(ctx, standard)
}
