package repositories

import (
	"ISO_Auditing_Tool/pkg/types"
	"context"
	"database/sql"
)

// DraftRepository is the concrete implementation
type RequirementRepository struct {
	db *sql.DB
}

// Ensure DraftRepository implements DraftRepositoryInterface
var _ RequirementRepositoryInterface = (*RequirementRepository)(nil)

func NewRequirementRepository(db *sql.DB) (RequirementRepositoryInterface, error) {
	return &RequirementRepository{db: db}, nil
}

func (r *RequirementRepository) GetByIDRequirement(ctx context.Context, requirement types.Requirement) (types.Requirement, error) {
	return types.Requirement{}, nil
}

func (r *RequirementRepository) GetByIDWithQuestionsRequirement(ctx context.Context, requirement types.Requirement) (types.Requirement, error) {
	return types.Requirement{}, nil
}

func (r *RequirementRepository) UpdateRequirement(ctx context.Context, requirement types.Requirement) (types.Requirement, error) {
	return types.Requirement{}, nil
}

func (r *RequirementRepository) UpdateRequirementAndDeleteDraft(ctx context.Context, requirement types.Requirement, draft types.Draft) (types.Requirement, error) {
	return types.Requirement{}, nil
}
